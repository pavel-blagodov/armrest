package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Request struct {
	url      string
	method   string
	username string
	password string
	query    url.Values
	ctx      context.Context
}

func request[T any](request Request) (T, error) {
	var RV T
	client := http.Client{}
	u, err := url.Parse(request.url)
	if err != nil {
		return RV, err
	}

	// if it's a GET, we need to append the query parameters.
	if request.method == "GET" {
		q := u.Query()
		for k, v := range request.query {
			q.Set(k, strings.Join(v, ","))
		}
		u.RawQuery = q.Encode()
	}

	// Define request
	req, err := http.NewRequestWithContext(request.ctx, request.method, u.String(), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating long-polling request: %v\n", err)
		return RV, err
	}
	req.SetBasicAuth(request.username, request.password)

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making long-polling request: %v\n", err)
		return RV, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return RV, err
	}

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Unexpected status code %d\n", resp.StatusCode)
		return RV, err
	}

	// Unmarshal JSON data into a struct
	json.Unmarshal(body, &RV)
	if err != nil {
		fmt.Println("Error Unmarshal data", err)
	}

	return RV, err
}
