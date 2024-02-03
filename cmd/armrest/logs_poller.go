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

func LogsPoller(flags *rootFlags) (func(cs ...chan Logs), func(), error) {
	var cancel context.CancelFunc

	stop := func() {
		if cancel != nil {
			cancel()
		}
	}

	// Start a goroutine to handle long polling
	start := func(cs ...chan Logs) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		doPoll := func() {
			queryParams := url.Values{}
			queryParams.Add("limit", "100")

			resp, err := get[Logs](ctx, Request{
				base:     flags.cbServerAPI,
				path:     "/logs",
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
