package app

import (
	"context"
	"fmt"
	"net/url"
)

func poolsDefaultPoller(flags *rootFlags) (func(), func(), chan PoolsDefault, error) {

	// Set up a channel to receive response
	ch := make(chan PoolsDefault)
	var cancel context.CancelFunc

	stop := func() {
		cancel()
	}

	// Start a goroutine to handle long polling
	start := func() {
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

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			resp, err := request[PoolsDefault](Request{
				url:      parsedURL.String(),
				method:   "GET",
				username: flags.username,
				password: flags.password,
				query:    queryParams,
				ctx:      ctx,
			})

			select {
			// Prepare next request
			case ch <- resp:
				etag = resp.Etag
			// Stop loop and request
			case <-ctx.Done():
				return
			}
		}
	}

	return start, stop, ch, nil
}
