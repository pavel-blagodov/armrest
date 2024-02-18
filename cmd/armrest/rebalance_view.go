package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/pavel-blagodov/armrest/cmd/utils"
)

func RebalanceWidget(progress int) (*donut.Donut, error) {
	var color cell.Color
	if progress > 0 {
		color = cell.ColorRed
	} else {
		progress = 100
		color = cell.ColorGray
	}
	widget, err := donut.New(
		donut.CellOpts(cell.FgColor(color)),
		donut.Label("Rebalance", cell.FgColor(cell.ColorGreen)),
	)
	if err := widget.Percent(progress); err != nil {
		panic(err)
	}
	return widget, err
}

func UpdateRebalanceLayout(ctx context.Context, t *tcell.Terminal, c *container.Container) (pdChannel chan TasksResponse) {
	ch := make(chan TasksResponse)

	go func() {
		for {
			select {
			case tasksResponse := <-ch:

				rebalanceTask := utils.Find[TasksItem](tasksResponse, func(task TasksItem) bool {
					return task.Type == "rebalance"
				})

				widget, err := RebalanceWidget(int(rebalanceTask.Progress))
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error new widget: %v\n", err)
				}
				textOptions := []container.Option{
					container.PlaceWidget(widget),
					container.BorderTitleAlignRight(),
				}

				if err := c.Update(rebalanceContainerID, textOptions...); err != nil {
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
