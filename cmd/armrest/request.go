package app

import (
	"bytes"
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
	base     string
	path     string
	method   string
	username string
	password string
	query    url.Values
}

func post[T any](ctx context.Context, request1 Request, payload []byte) (T, error) {
	request1.method = "POST"
	return request[T](ctx, request1, payload)
}

func get[T any](ctx context.Context, request1 Request) (T, error) {
	request1.method = "GET"
	return request[T](ctx, request1, nil)
}

func request[T any](ctx context.Context, request Request, payload []byte) (T, error) {
	var rv T
	client := http.Client{}
	u, err := url.Parse(request.base)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing base URL: %v\n", err)
		return rv, err
	}

	// if it's a GET, we need to append the query parameters.
	if request.method == "GET" {
		q := u.Query()
		for k, v := range request.query {
			q.Set(k, strings.Join(v, ","))
		}
		u.RawQuery = q.Encode()
	}

	u.Path = request.path

	// Define request
	req, err := http.NewRequestWithContext(ctx, request.method, u.String(), bytes.NewBuffer(payload))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating long-polling request: %v\n", err)
		return rv, err
	}
	req.SetBasicAuth(request.username, request.password)

	if request.method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making long-polling request: %v\n", err)
		return rv, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return rv, err
	}

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Unexpected status code %d\n", resp.StatusCode)
		return rv, err
	}

	// Unmarshal JSON data into a struct
	json.Unmarshal(body, &rv)
	if err != nil {
		fmt.Println("Error Unmarshal data", err)
	}

	return rv, err
}
