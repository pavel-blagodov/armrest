package app

import (
	"context"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/pavel-blagodov/armrest/cmd/utils"
)

func UpdateKvOpsLayout(ctx context.Context, t *tcell.Terminal, c *container.Container) (pdChannel chan RangeResponse) {
	ch := make(chan RangeResponse)

	widget, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorCyan)),
	)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case rangeResponse := <-ch:

				series := utils.MapSlice[[]any, float64](rangeResponse[0].Data[0].Values, func(value []any) float64 {
					strValue, ok := value[1].(string)
					if !ok {
						return math.NaN()
					}
					floatValue, err := strconv.ParseFloat(strValue, 64)
					if err != nil {
						return math.NaN()
					}
					return floatValue
				})

				times := utils.Reduce[[]any, map[int]string](rangeResponse[0].Data[0].Values, func(acc map[int]string, value []any, index int) map[int]string {
					ts, ok := value[0].(float64)
					if !ok {
						acc[index] = "-"
						return acc
					}
					t := time.Unix(int64(ts), 0)

					acc[index] = t.Format("3:04 PM")
					return acc
				}, map[int]string{})

				if err := widget.Series("first", series,
					linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(33))),
					linechart.SeriesXLabels(times),
				); err != nil {
					panic(err)
				}

				textOptions := []container.Option{
					container.PlaceWidget(widget),
					container.BorderTitleAlignRight(),
				}

				if err := c.Update(kvOpsContainerID, textOptions...); err != nil {
					fmt.Fprintf(os.Stderr, "Error update: %v\n", err)
				}

			case <-ctx.Done():
				return

			case <-time.After(10 * time.Minute):
				fmt.Println("No notifications received for 10 minutes. Exiting.")
				return
			}
		}
	}()

	return ch
}
