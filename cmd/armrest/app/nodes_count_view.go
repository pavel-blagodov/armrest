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
)

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

func updateNodesServiceCountLayout(ctx context.Context, t *tcell.Terminal, c *container.Container, pdChannel chan PoolsDefault) {
	//nodes system stats
	for {
		select {
		case poolsDefault := <-pdChannel:
			nodesCountWidgets, err := newNodesServiceCountWidget(ctx, poolsDefault)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error new widget: %v\n", err)
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
