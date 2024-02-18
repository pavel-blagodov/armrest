package app

import (
	"context"
	"time"
)

type RangeResponse []struct {
	Data []struct {
		Metric struct {
			Nodes []string `json:"nodes"`
		} `json:"metric"`
		Values [][]any `json:"values"`
	} `json:"data"`
	Errors         []interface{} `json:"errors"`
	StartTimestamp int           `json:"startTimestamp"`
	EndTimestamp   int           `json:"endTimestamp"`
}

// type RangeRequest struct {
// 	Step             int                  `json:"step"`
// 	TimeWindow       string               `json:"timeWindow"`
// 	Start            int                  `json:"start"`
// 	Metric           []RangeRequestMetric `json:"metric"`
// 	NodesAggregation string               `json:"nodesAggregation"`
// 	ApplyFunctions   []string             `json:"applyFunctions,omitempty"`
// 	AlignTimestamps  bool                 `json:"alignTimestamps"`
// }
// type RangeRequestMetric struct {
// 	Label string `json:"label"`
// 	Value string `json:"value"`
// }

func RangePoller(flags *rootFlags, payload []byte) (func(cs ...chan RangeResponse), func(), error) {
	var cancel context.CancelFunc

	stop := func() {
		if cancel != nil {
			cancel()
		}
	}

	// Start a goroutine to handle long polling
	start := func(cs ...chan RangeResponse) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		doPoll := func() {

			resp, err := post[RangeResponse](ctx, Request{
				base:     flags.cbServerAPI,
				path:     "/pools/default/stats/range/",
				username: flags.username,
				password: flags.password,
			}, payload)

			if err != nil {
				return
			}
			for _, c := range cs {
				select {
				case c <- resp:
				}
			}
		}
		ticker := time.NewTicker(5 * time.Second)
		defer func() { ticker.Stop() }()
		doPoll()

		for {
			select {
			case <-ticker.C:
				doPoll()
			case <-ctx.Done():
				return
			}
		}
	}

	return start, stop, nil
}
