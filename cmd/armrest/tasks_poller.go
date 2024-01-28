package app

import (
	"context"
	"math"
	"time"
)

func findRecommendedRefreshPeriod(tasks TasksResponse) float64 {
	if len(tasks) == 0 {
		return 5.0
	}

	minRefreshPeriod := math.MaxFloat64

	for _, obj := range tasks {
		if obj.RecommendedRefreshPeriod < minRefreshPeriod {
			minRefreshPeriod = obj.RecommendedRefreshPeriod
		}
	}
	if minRefreshPeriod == 0 {
		return 5.0
	}
	return minRefreshPeriod
}

func TasksPoller(flags *rootFlags) (func(cs ...chan TasksResponse), func(), chan PoolsDefault, error) {
	var cancel context.CancelFunc
	poolsDefaultCh := make(chan PoolsDefault)

	stop := func() {
		if cancel != nil {
			cancel()
		}
	}

	// Start a goroutine to handle long polling
	start := func(cs ...chan TasksResponse) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ticker := time.NewTicker(1 * time.Second)

		doPoll := func() {

			resp, err := request[TasksResponse](ctx, Request{
				method:   "GET",
				base:     flags.cbServerAPI,
				path:     "/pools/default/tasks",
				username: flags.username,
				password: flags.password,
			})

			if err != nil {
				return
			}

			ticker.Reset(time.Duration(findRecommendedRefreshPeriod(resp) * float64(time.Second)))

			for _, c := range cs {
				select {
				case c <- resp:

				}
			}
		}

		defer func() { ticker.Stop() }()

		var url = ""

		for {
			select {
			case poolDefault := <-poolsDefaultCh:
				if url != poolDefault.Tasks.URI {
					doPoll()
					url = poolDefault.Tasks.URI
				}
			case <-ticker.C:
				if url != "" {
					doPoll()
				}
			case <-ctx.Done():
				return
			}
		}
	}

	return start, stop, poolsDefaultCh, nil
}
