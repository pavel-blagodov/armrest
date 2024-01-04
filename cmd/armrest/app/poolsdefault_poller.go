package app

import (
	"context"
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
		var etag string
		for {
			queryParams := url.Values{}
			queryParams.Add("etag", etag)
			queryParams.Add("waitChange", "10000")

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			resp, err := request[PoolsDefault](ctx, Request{
				method:   "GET",
				base:     flags.cbServerAPI,
				path:     "/pools/default",
				query:    queryParams,
				username: flags.username,
				password: flags.password,
			})

			if err != nil {
				return
			}

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
