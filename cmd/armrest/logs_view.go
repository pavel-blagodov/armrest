package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/widgets/text"

	"github.com/pavel-blagodov/armrest/cmd/utils"
)

func newLogsWidget(ctx context.Context, logs Logs) (*text.Text, error) {
	wrapped, err := text.New(text.WrapAtRunes())
	formattedTime := "15:04 02 Jan"

	for i := range logs.List {
		log := logs.List[i]
		parsedTime, err := time.Parse(time.RFC3339Nano, log.ServerTime)
		if err != nil {
			fmt.Println("Error parsing timestamp:", err)
			return wrapped, err
		}
		if err := wrapped.Write(fmt.Sprintf("%s %s%s", parsedTime.Format(formattedTime), log.Module, "\n"), text.WriteCellOpts(cell.FgColor(cell.ColorMagenta))); err != nil {
			return nil, err
		}

		if err := wrapped.Write(utils.TrimString(log.Text, 200, "...") + "\n\n"); err != nil {
			return nil, err
		}
	}

	return wrapped, err
}

func UpdateLogsLayout(ctx context.Context, t *tcell.Terminal, c *container.Container) (logsChannel chan Logs) {
	ch := make(chan Logs)

	go func() {
		for {
			select {
			case logs := <-ch:
				logsWidgets, err := newLogsWidget(ctx, logs)

				if err != nil {
					fmt.Fprintf(os.Stderr, "Error new widget: %v\n", err)
				}

				textOptions := []container.Option{
					container.Border(linestyle.Light),
					container.BorderTitle("Logs"),
					container.PlaceWidget(logsWidgets),
					container.BorderTitleAlignRight(),
				}

				if err := c.Update(logsContainerID, textOptions...); err != nil {
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

	//render straight away
	ch <- Logs{
		List: []Log{},
	}

	return ch
}
