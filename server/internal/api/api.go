package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// getClient returns a configured *http.Client with a default timeout of 5 seconds.
func getClient() *http.Client {
	return &http.Client{Timeout: 5 * time.Second}
}

// GetDataFromApi sends an HTTP GET request to the specified URL and returns the response body as a byte slice or an error.
func GetDataFromApi(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	res, err := getClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		snippet, _ := io.ReadAll(io.LimitReader(res.Body, 8<<10))
		return nil, fmt.Errorf("unexpected status %d: %s", res.StatusCode, string(snippet))
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return b, nil
}
