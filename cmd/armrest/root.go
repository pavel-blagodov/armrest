package app

import (
	"context"
	"fmt"
	"os"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/button"
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

func rootCmd(flags *rootFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		t, err := tcell.New(tcell.ColorMode(terminalapi.ColorMode256))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error tcell: %v\n", err)
		}
		defer t.Close()

		ctx, cancel := context.WithCancel(context.Background())

		quitter := func(k *terminalapi.Keyboard) {
			if k.Key == 'q' || k.Key == 'Q' {
				cancel()
			}
		}

		rootContainer, err := crateGeneralLayout(t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error new container: %v\n", err)
		}

		buttons, err := newLayoutButtons(rootContainer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating buttons: %v\n", err)
		}

		poolDefaultStart, _, err := poolsDefaultPoller(flags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error performing authenticated long-polling: %v\n", err)
		}
		logsStart, _, err := logsPoller(flags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error performing authenticated long-polling: %v\n", err)
		}

		updateButtonsLayout(buttons, rootContainer)

		logsCh := updateLogsLayout(ctx, t, rootContainer)
		nodesCountCh := updateNodesServiceCountLayout(ctx, t, rootContainer)
		nodesCh := updateNodesLayout(ctx, t, rootContainer)

		go poolDefaultStart(nodesCountCh, nodesCh)

		go logsStart(logsCh)

		if err := termdash.Run(ctx, t, rootContainer, termdash.KeyboardSubscriber(quitter)); err != nil {
			fmt.Fprintf(os.Stderr, "Error termdash.Run: %v\n", err)
		}

		// poolsResp, err := request[Pools](ctx, Request{
		// 	base:     flags.cbServerAPI,
		// 	path:     "/pools",
		// 	method:   "GET",
		// 	username: flags.username,
		// 	password: flags.password,
		// })
	}
}

// layoutType represents the possible layouts the buttons switch between.
type layoutType int

const (
	// layoutAll displays all the widgets.
	layoutServer layoutType = iota
	// layoutText focuses onto the text widget.
	layoutBuckets
	// layoutSparkLines focuses onto the sparklines.
	layoutXDCR
)

const rootContainerID = "root"
const nodesSystemStatsContainerID = "nodesSystemStats"
const nodesCountContainerID = "nodesCountContainerID"
const layoutSpecificContainerID = "layoutSpecificContainerID"
const layoutButtonsContainerID = "layoutButtonsContainerID"
const logsContainerID = "logsContainerID"

func crateGeneralLayout(t *tcell.Terminal) (*container.Container, error) {
	rootContainer, err := container.New(t,
		container.ID(rootContainerID),
		container.SplitVertical(
			container.Left(
				container.SplitHorizontal(
					container.Top(
						container.Border(linestyle.Light),
						container.BorderTitle("Press Q to quit"),
					),
					container.Bottom(
						container.SplitHorizontal(
							container.Top(container.ID(layoutButtonsContainerID)),
							container.Bottom(container.ID(layoutSpecificContainerID)),
							container.SplitPercent(20),
						),
					),
					container.SplitPercent(30),
				),
			),
			container.Right(
				container.SplitHorizontal(
					container.Top(container.ID(nodesCountContainerID)),
					container.Bottom(container.ID(logsContainerID)),
					container.SplitPercent(30),
				),
			),
			container.SplitPercent(70),
		),
	)
	return rootContainer, err
}

// setLayout sets the specified layout.
func setLayout(c *container.Container, lt layoutType) error {
	switch lt {
	case layoutServer:
	case layoutBuckets:
	case layoutXDCR:
	}
	return nil
}

// layoutButtons are buttons that change the layout.
type layoutButtons struct {
	serversB *button.Button
	bucketsB *button.Button
	xdcrB    *button.Button
}

func updateButtonsLayout(buttons *layoutButtons, c *container.Container) error {
	builder := grid.New()
	builder.Add(grid.RowHeightPerc(5,
		grid.ColWidthPerc(33,
			grid.Widget(buttons.serversB),
		),
		grid.ColWidthPerc(33,
			grid.Widget(buttons.bucketsB),
		),
		grid.ColWidthPerc(33,
			grid.Widget(buttons.xdcrB),
		),
	))
	gridOpts, err := builder.Build()
	if err != nil {
		return err
	}
	if err := c.Update(layoutButtonsContainerID, gridOpts...); err != nil {
		return err
	}
	return nil
}

// newLayoutButtons returns buttons that dynamically switch the layouts.
func newLayoutButtons(c *container.Container) (*layoutButtons, error) {
	opts := []button.Option{
		button.WidthFor("Servers"),
		button.FillColor(cell.ColorNumber(220)),
		button.Height(1),
	}

	serversB, err := button.New("Servers", func() error {
		return setLayout(c, layoutServer)
	}, opts...)
	if err != nil {
		return nil, err
	}

	bucketsB, err := button.New("Buckets", func() error {
		return setLayout(c, layoutBuckets)
	}, opts...)
	if err != nil {
		return nil, err
	}

	xdcrB, err := button.New("XDCR", func() error {
		return setLayout(c, layoutXDCR)
	}, opts...)
	if err != nil {
		return nil, err
	}

	return &layoutButtons{
		serversB: serversB,
		bucketsB: bucketsB,
		xdcrB:    xdcrB,
	}, nil
}
