package app

import (
	"context"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/barchart"
	"github.com/spf13/cobra"
)

type rootFlags struct {
	username    string
	password    string
	cbServerAPI string
}

func NewRootCommand() *cobra.Command {
	flags := rootFlags{}
	cmd := &cobra.Command{
		Use:   "armrest",
		Short: "Superduper couchbase server CLI dashboard",
		Long:  "Superduper couchbase server CLI dashboard",
		Run:   rootCmd(&flags),
	}

	cmd.PersistentFlags().StringVar(&flags.username, "username", "", "Username for authentication")
	cmd.PersistentFlags().StringVar(&flags.password, "password", "", "Password for authentication")
	cmd.PersistentFlags().StringVar(&flags.cbServerAPI, "url", "", "Couchbase server http API URL")

	return cmd
}

// redrawInterval is how often termdash redraws the screen.
const redrawInterval = 250 * time.Millisecond
const rootID = "root"

func rootCmd(flags *rootFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		pdStart, pdChannel, err := poolsDefaultPoller(flags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error performing authenticated long-polling: %v\n", err)
		}

		t, err := tcell.New(tcell.ColorMode(terminalapi.ColorMode256))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error tcell: %v\n", err)
		}
		defer t.Close()

		c, err := container.New(t, container.ID(rootID))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error new container: %v\n", err)
		}

		ctx, cancel := context.WithCancel(context.Background())

		quitter := func(k *terminalapi.Keyboard) {
			if k.Key == 'q' || k.Key == 'Q' {
				cancel()
			}
		}

		go func() {
			for {
				select {
				case poolsDefault := <-pdChannel:

					w, err := newWidgets(ctx, t, c, poolsDefault)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error new widget: %v\n", err)
					}

					gridOpts, err := gridLayout(w, poolsDefault) // equivalent to contLayout(w)

					for i, node := range poolsDefault.Nodes {
						values := []int{getCpuUsage(node), getRamUsage(node), getSwapUsage(node)}
						w.barCharts[i].Values(values, 100)
					}

					if err := c.Update(rootID, gridOpts...); err != nil {
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

		go pdStart()

		if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter)); err != nil {
			fmt.Fprintf(os.Stderr, "Error termdash.Run: %v\n", err)
		}
	}
}

func getRamUsage(node Nodes) int {
	total := node.MemoryTotal
	free := node.MemoryFree
	used := total - free

	if total == 0 || math.IsInf(float64(free), 0) || math.IsNaN(float64(free)) {
		return 0
	}
	return int(float64(used) / float64(total) * 100)

}
func getSwapUsage(node Nodes) int {
	total := node.SystemStats.SwapTotal
	used := node.SystemStats.SwapUsed

	if total == 0 || math.IsInf(float64(used), 0) || math.IsNaN(float64(used)) {
		return 0
	}
	return int(float64(used) / float64(total) * 100)
}
func getCpuUsage(node Nodes) int {
	var cpuRate = node.SystemStats.CPUUtilizationRate
	return int(math.Floor(cpuRate*100) / 100)
}

// widgets holds the widgets used by this demo.
type widgets struct {
	barCharts []*barchart.BarChart
}

func newWidgets(ctx context.Context, t terminalapi.Terminal, c *container.Container, poolDefault PoolsDefault) (*widgets, error) {
	bc, err := newBarChart(ctx, poolDefault)
	if err != nil {
		return nil, err
	}

	return &widgets{
		barCharts: bc,
	}, nil
}

// newBarChart returns a BarcChart that displays random values on multiple bars.
func newBarChart(ctx context.Context, poolsDefault PoolsDefault) ([]*barchart.BarChart, error) {
	rv := make([]*barchart.BarChart, len(poolsDefault.Nodes))

	for i, _ := range poolsDefault.Nodes {
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

func gridLayout(w *widgets, poolsDefault PoolsDefault) ([]container.Option, error) {
	leftRows := make([]grid.Element, len(poolsDefault.Nodes))

	builder := grid.New()
	length := len(poolsDefault.Nodes)
	for i, node := range poolsDefault.Nodes {
		leftRows[i] = grid.RowHeightPerc(100/length,
			grid.Widget(w.barCharts[i],
				container.Border(linestyle.Light),
				container.BorderTitle(node.Hostname),
				container.BorderTitleAlignRight(),
			))
	}

	builder.Add(
		grid.ColWidthPerc(30, leftRows...),
	)

	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return gridOpts, nil
}
