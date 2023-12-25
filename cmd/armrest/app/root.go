package app

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

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

		pdStart, pdch, err := poolsDefaultPoller(flags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error performing authenticated long-polling: %v\n", err)
		}
		go pdStart()

		// Wait for notifications and print them
		// go func() {
		for {
			select {
			case notification := <-pdch:
				response, _ := json.Marshal(notification)
				fmt.Println("Received notification:", string(response))

				// TODO: Add your logic to handle the received notification
				// You may want to update the main loop or perform other actions here

			case <-time.After(10 * time.Minute):
				fmt.Println("No notifications received for 10 minutes. Exiting.")
				return
			}
		}
		// }()
	}
}
