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
	"github.com/mum4k/termdash/widgets/text"
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
		grid.ColWidthPerc(50, leftRows...),
	)

	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return gridOpts, nil
}

func newNodesServiceCountWidget(ctx context.Context, poolsDefault PoolsDefault) (*text.Text, error) {
	serverIcon := "\u2630"
	//nodes count
	wrapped, err := text.New(text.WrapAtRunes())
	if err != nil {
		return nil, err
	}
	//active
	if err := wrapped.Write(serverIcon, text.WriteCellOpts(cell.FgColor(cell.ColorGreen))); err != nil {
		return nil, err
	}
	failedOver := filter[Node](poolsDefault.Nodes, func(node Node) bool {
		return node.ClusterMembership == "inactiveFailed"
	})
	onlyActive := filter[Node](poolsDefault.Nodes, func(node Node) bool {
		return node.ClusterMembership == "active"
	})
	active := append(failedOver, onlyActive...)

	if err := wrapped.Write(fmt.Sprintf(" %d %s %s%s", len(active), "active", pluralize(len(active), "node", "nodes"), "\n")); err != nil {
		return nil, err
	}

	//failed over
	if err := wrapped.Write(serverIcon, text.WriteCellOpts(cell.FgColor(cell.ColorYellow))); err != nil {
		return nil, err
	}
	if err := wrapped.Write(fmt.Sprintf(" %d %s %s%s", len(failedOver), "failed-over", pluralize(len(failedOver), "node", "nodes"), "\n")); err != nil {
		return nil, err
	}

	//pending
	pending := filter[Node](poolsDefault.Nodes, func(node Node) bool {
		return node.ClusterMembership != "active"
	})
	if err := wrapped.Write(serverIcon, text.WriteCellOpts(cell.FgColor(cell.ColorYellow))); err != nil {
		return nil, err
	}
	if err := wrapped.Write(fmt.Sprintf(" %d %s %s", len(pending), pluralize(len(pending), "node", "nodes"), "pending rebalance\n")); err != nil {
		return nil, err
	}

	//down
	down := filter[Node](poolsDefault.Nodes, func(node Node) bool {
		return node.Status != "healthy"
	})
	if err := wrapped.Write(serverIcon, text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		return nil, err
	}
	if err := wrapped.Write(fmt.Sprintf(" %d %s %s%s", len(down), "inactive", pluralize(len(down), "node", "nodes"), "\n")); err != nil {
		return nil, err
	}

	return wrapped, nil
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

			//nodes system stats
			nodesCountWidgets, err := newNodesServiceCountWidget(ctx, poolsDefault)
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

			if err := c.Update(nodesSystemStatsContainerID, gridOpts...); err != nil {
				fmt.Fprintf(os.Stderr, "Error update: %v\n", err)
			}

			//nodes count
			textOptions := []container.Option{
				container.Border(linestyle.Light),
				container.BorderTitle("Nodes status"),
				container.PlaceWidget(nodesCountWidgets),
			}
			if err := c.Update(nodesCountContainerID, textOptions...); err != nil {
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
