package app

import (
	"context"
	"net/url"
	"time"
)

type Log struct {
	Node       string `json:"node"`
	Type       string `json:"type"`
	Code       int    `json:"code"`
	Module     string `json:"module"`
	Tstamp     int64  `json:"tstamp"`
	ShortText  string `json:"shortText"`
	Text       string `json:"text"`
	ServerTime string `json:"serverTime"`
}

type Logs struct {
	List []Log `json:"list"`
}

func logsPoller(flags *rootFlags) (func(), func(), chan Logs, error) {
	// Set up a channel to receive response
	ch := make(chan Logs)
	var cancel context.CancelFunc

	stop := func() {
		cancel()
	}

	// Start a goroutine to handle long polling
	start := func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		doPoll := func() {
			queryParams := url.Values{}
			queryParams.Add("limit", "100")

			resp, err := request[Logs](ctx, Request{
				method:   "GET",
				base:     flags.cbServerAPI,
				path:     "/logs",
				query:    queryParams,
				username: flags.username,
				password: flags.password,
			})

			if err != nil {
				return
			}
			ch <- resp
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

	return start, stop, ch, nil
}
