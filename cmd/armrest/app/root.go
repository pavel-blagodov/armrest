package app

import (
	"context"
	"fmt"
	"os"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
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

const rootContainerID = "root"
const nodesContainerID = "nodes"

func rootCmd(flags *rootFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		poolDefaultStart, _, poolDefaultChannel, err := poolsDefaultPoller(flags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error performing authenticated long-polling: %v\n", err)
		}

		t, err := tcell.New(tcell.ColorMode(terminalapi.ColorMode256))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error tcell: %v\n", err)
		}
		defer t.Close()

		rootContainer, err := container.New(t,
			container.ID(rootContainerID),
			container.SplitVertical(
				container.Left(container.ID(nodesContainerID)),
				container.Right(),
			),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error new container: %v\n", err)
		}

		ctx, cancel := context.WithCancel(context.Background())

		quitter := func(k *terminalapi.Keyboard) {
			if k.Key == 'q' || k.Key == 'Q' {
				cancel()
			}
		}

		go poolDefaultStart()

		go updateNodesLayout(ctx, t, rootContainer, poolDefaultChannel)

		if err := termdash.Run(ctx, t, rootContainer, termdash.KeyboardSubscriber(quitter)); err != nil {
			fmt.Fprintf(os.Stderr, "Error termdash.Run: %v\n", err)
		}
	}
}
