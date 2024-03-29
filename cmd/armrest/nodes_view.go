package app

import (
	"context"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/widgets/barchart"
)

func clearLayout() ([]container.Option, error) {
	builder := grid.New()
	leftRows := make([]grid.Element, 1)
	leftRows[0] = grid.RowHeightPerc(99,
		grid.Widget(nil,
			container.Border(linestyle.None),
			container.BorderTitle(""),
			container.BorderTitleAlignRight(),
		))
	builder.Add(grid.ColWidthPerc(99, leftRows...))
	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return gridOpts, nil
}

func newNodesLayout(widget []*barchart.BarChart, poolsDefault PoolsDefault) ([]container.Option, error) {
	leftRows := make([]grid.Element, len(poolsDefault.Nodes))

	builder := grid.New()
	length := len(poolsDefault.Nodes)
	for i, node := range poolsDefault.Nodes {
		leftRows[i] = grid.RowHeightPerc(99/length,
			grid.Widget(widget[i],
				container.Border(linestyle.Light),
				container.BorderTitle(node.Hostname),
				container.BorderTitleAlignRight(),
			))
	}

	builder.Add(
		grid.ColWidthPerc(99, leftRows...),
	)

	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return gridOpts, nil
}

func newNodesSystemsStatsWidgets(ctx context.Context, poolsDefault PoolsDefault) ([]*barchart.BarChart, error) {
	rv := make([]*barchart.BarChart, len(poolsDefault.Nodes))

	for i := range poolsDefault.Nodes {
		bc, err := barchart.New(
			barchart.BarColors([]cell.Color{
				cell.ColorNumber(124),
				cell.ColorNumber(174),
				cell.ColorNumber(194),
			}),
			barchart.ValueColors([]cell.Color{
				cell.ColorWhite,
				cell.ColorWhite,
				cell.ColorWhite,
			}),
			barchart.ShowValues(),
			barchart.BarWidth(0),
			barchart.Labels([]string{
				"CPU",
				"RAM",
				"SWAP",
			}),
		)
		if err != nil {
			return nil, err
		}

		rv[i] = bc
	}

	return rv, nil
}

func UpdateNodesLayout(ctx context.Context, t *tcell.Terminal, c *container.Container) (pdChannel chan PoolsDefault) {
	ch := make(chan PoolsDefault)

	go func() {
		for {
			select {
			case poolsDefault := <-pdChannel:

				nodesSystemsStatsWidgets, err := newNodesSystemsStatsWidgets(ctx, poolsDefault)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error new widget: %v\n", err)
				}

				clearGridOpts, err := clearLayout() // equivalent to contLayout(w)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error new layout: %v\n", err)
				}

				if err := c.Update(layoutSpecificContainerID, clearGridOpts...); err != nil {
					fmt.Fprintf(os.Stderr, "Error update: %v\n", err)
				}

				gridOpts, err := newNodesLayout(nodesSystemsStatsWidgets, poolsDefault) // equivalent to contLayout(w)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error new layout: %v\n", err)
				}

				for i, node := range poolsDefault.Nodes {
					values := []int{getCpuUsage(node), getRamUsage(node), getSwapUsage(node)}
					nodesSystemsStatsWidgets[i].Values(values, 100)
				}

				if err := c.Update(layoutSpecificContainerID, gridOpts...); err != nil {
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

func getRamUsage(node Node) int {
	total := node.MemoryTotal
	free := node.MemoryFree
	used := total - free

	if total == 0 || math.IsInf(float64(free), 0) || math.IsNaN(float64(free)) {
		return 0
	}
	return int(float64(used) / float64(total) * 100)
}

func getSwapUsage(node Node) int {
	total := node.SystemStats.SwapTotal
	used := node.SystemStats.SwapUsed

	if total == 0 || math.IsInf(float64(used), 0) || math.IsNaN(float64(used)) {
		return 0
	}
	return int(float64(used) / float64(total) * 100)
}

func getCpuUsage(node Node) int {
	var cpuRate = node.SystemStats.CPUUtilizationRate
	return int(math.Floor(cpuRate*100) / 100)
}
