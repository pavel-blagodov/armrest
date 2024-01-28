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

func writeToWrapper(wrapped *text.Text, length int, color cell.Color, kind string) error {
	if err := wrapped.Write("\u2630", text.WriteCellOpts(cell.FgColor(color))); err != nil {
		return err
	}
	if err := wrapped.Write(fmt.Sprintf(" %d %s %s%s", length, kind, utils.Pluralize(length, "node", "nodes"), "\n")); err != nil {
		return err
	}
	return nil
}

func newNodesServiceCountWidget(ctx context.Context, poolsDefault PoolsDefault) (*text.Text, error) {
	//nodes count
	wrapped, err := text.New(text.WrapAtRunes())
	if err != nil {
		return nil, err
	}
	//active
	failedOver := utils.Filter[Node](poolsDefault.Nodes, func(node Node) bool {
		return node.ClusterMembership == "inactiveFailed"
	})
	onlyActive := utils.Filter[Node](poolsDefault.Nodes, func(node Node) bool {
		return node.ClusterMembership == "active"
	})
	active := append(failedOver, onlyActive...)
	//pending
	pending := utils.Filter[Node](poolsDefault.Nodes, func(node Node) bool {
		return node.ClusterMembership != "active"
	})
	//down
	down := utils.Filter[Node](poolsDefault.Nodes, func(node Node) bool {
		return node.Status != "healthy"
	})

	if err := writeToWrapper(wrapped, len(active), cell.ColorGreen, "active"); err != nil {
		return nil, err
	}
	if err := writeToWrapper(wrapped, len(failedOver), cell.ColorYellow, "failed-over"); err != nil {
		return nil, err
	}
	if err := writeToWrapper(wrapped, len(pending), cell.ColorYellow, "pending rebalance"); err != nil {
		return nil, err
	}
	if err := writeToWrapper(wrapped, len(down), cell.ColorRed, "inactive"); err != nil {
		return nil, err
	}

	return wrapped, nil
}

func UpdateNodesServiceCountLayout(ctx context.Context, t *tcell.Terminal, c *container.Container) (pdChannel chan PoolsDefault) {
	ch := make(chan PoolsDefault)

	go func() {
		for {
			select {
			case poolsDefault := <-ch:
				nodesCountWidgets, err := newNodesServiceCountWidget(ctx, poolsDefault)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error new widget: %v\n", err)
				}

				textOptions := []container.Option{
					container.Border(linestyle.Light),
					container.BorderTitle("Nodes status"),
					container.PlaceWidget(nodesCountWidgets),
					container.BorderTitleAlignRight(),
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
	}()

	ch <- PoolsDefault{
		Nodes: []Node{},
	}

	return ch
}
