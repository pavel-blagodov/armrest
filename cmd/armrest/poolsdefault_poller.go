package app

import (
	"context"
	"net/url"
)

func PoolsDefaultPoller(flags *rootFlags) (func(cs ...chan PoolsDefault), func(), error) {
	var cancel context.CancelFunc

	stop := func() {
		if cancel != nil {
			cancel()
		}
	}

	// Start a goroutine to handle long polling
	start := func(cs ...chan PoolsDefault) {
		var etag string
		for {
			queryParams := url.Values{}
			queryParams.Add("etag", etag)
			queryParams.Add("waitChange", "10000")

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			resp, err := get[PoolsDefault](ctx, Request{
				base:     flags.cbServerAPI,
				path:     "/pools/default",
				query:    queryParams,
				username: flags.username,
				password: flags.password,
			})

			if err != nil {
				return
			}

			for _, c := range cs {
				select {
				case c <- resp:
				case <-ctx.Done():
					return
				}
			}
			etag = resp.Etag
		}
	}

	return start, stop, nil
}
