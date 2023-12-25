package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func poolsDefaultPoller(flags *rootFlags) (func(), chan PoolsDefault, error) {
	client := http.DefaultClient

	// Set up a channel to receive response
	ch := make(chan PoolsDefault)

	// Start a goroutine to handle long polling
	start := func() {
		fmt.Println("start")
		var etag = ""
		for {
			// Construct url
			parsedURL, err := url.Parse(flags.cbServerAPI)
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

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Define request
			req, err := http.NewRequestWithContext(ctx, "GET", longPollUrl, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating long-polling request: %v\n", err)
				return
			}
			req.SetBasicAuth(flags.username, flags.password)

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
			var poolsDefault PoolsDefault
			json.Unmarshal(body, &poolsDefault)
			if err != nil {
				fmt.Println("Error Unmarshal data", err)
			}

			select {
			// Prepare next request
			case ch <- poolsDefault:
				fmt.Println("poll")
				etag = poolsDefault.Etag
			// Stop loop and request
			case <-ch:
				fmt.Println("cancel")
				cancel()
				return
			}
		}
	}

	return start, ch, nil
}
