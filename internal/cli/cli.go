package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	types "github.com/pavel-blagodov/armrest/internal/http/poolsdefault"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "armrest",
	Short: "My CLI tool",
	Long:  "A simple CLI tool written in Go",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from mycli!")

		err := performAuthenticatedLongPolling()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error performing authenticated long-polling: %v\n", err)
		}
	},
}

var username string
var password string
var httpApiUrl string

func init() {
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "Username for authentication")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Password for authentication")
	rootCmd.PersistentFlags().StringVar(&httpApiUrl, "url", "", "Couchbase server http/https url")
}

// Execute is the entry point for the CLI tool
func Execute() error {
	return rootCmd.Execute()
}

func performAuthenticatedLongPolling() error {
	client := http.DefaultClient

	// Set up a channel to receive response
	notificationCh := make(chan types.PoolsDefault)
	var etag = ""

	// Start a goroutine to handle long polling
	go func() {
		for {
			//Construct url
			parsedURL, err := url.Parse(httpApiUrl)
			if err != nil {
				fmt.Println("Error parsing URL:", err)
				return
			}
			parsedURL.Path = "/pools/default"
			queryParams := url.Values{}
			queryParams.Add("etag", etag)
			queryParams.Add("waitChange", "10000")
			parsedURL.RawQuery = queryParams.Encode()
			var longPollUrl = parsedURL.String()

			fmt.Println("LongPollUrl:", longPollUrl)

			//Define request
			req, err := http.NewRequest("GET", longPollUrl, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating long-polling request: %v\n", err)
				return
			}
			req.SetBasicAuth(username, password)

			// Make request
			resp, err := client.Do(req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error making long-polling request: %v\n", err)
				return
			}
			defer resp.Body.Close()

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}

			// Check the HTTP status code
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Error: Unexpected status code %d\n", resp.StatusCode)
				return
			}

			// Unmarshal JSON data into a struct
			var poolsDefault types.PoolsDefault
			json.Unmarshal(body, &poolsDefault)
			if err != nil {
				fmt.Println("error:", err)
			}

			//update next etag
			etag = poolsDefault.Etag

			notificationCh <- poolsDefault
		}
	}()

	// Wait for notifications and print them
	for {
		select {
		case notification := <-notificationCh:
			response, _ := json.Marshal(notification)
			fmt.Println("Received notification:", string(response))

			// TODO: Add your logic to handle the received notification
			// You may want to update the main loop or perform other actions here

		case <-time.After(10 * time.Minute):
			fmt.Println("No notifications received for 10 minutes. Exiting.")
			return nil
		}
	}
}
