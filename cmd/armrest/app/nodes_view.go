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

func newNodesLayout(widget []*barchart.BarChart, poolsDefault PoolsDefault) ([]container.Option, error) {
	leftRows := make([]grid.Element, len(poolsDefault.Nodes))

	builder := grid.New()
	length := len(poolsDefault.Nodes)
	for i, node := range poolsDefault.Nodes {
		leftRows[i] = grid.RowHeightPerc(100/length,
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

func updateNodesLayout(ctx context.Context, t *tcell.Terminal, c *container.Container, pdChannel chan PoolsDefault) {
	for {
		select {
		case poolsDefault := <-pdChannel:

			//nodes system stats
			nodesSystemsStatsWidgets, err := newNodesSystemsStatsWidgets(ctx, poolsDefault)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error new widget: %v\n", err)
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
