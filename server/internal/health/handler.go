package health

import (
	"aggregator/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HealthHandler handles the health check endpoint by verifying the availability of two external servers.
// Responds with JSON indicating the health status based on the servers' availability.
func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("[GET] /health", time.Now().Format("2006-01-02 15:04:05"))

	// Configure HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second, // Global timeout of 5 seconds
	}

	// Check SERVER1
	if !checkServer(client, config.SERVER1_URL) {
		err := json.NewEncoder(w).Encode(NOk())
		if err != nil {
			return
		}
		return
	}

	// Check SERVER2
	if !checkServer(client, config.SERVER2_URL) {
		err := json.NewEncoder(w).Encode(NOk())
		if err != nil {
			return
		}
		return
	}

	// If we reach here, both servers are OK
	err := json.NewEncoder(w).Encode(Ok())
	if err != nil {
		return
	}
}

// checkServer sends a GET request to the specified URL using the provided HTTP client and checks if the server is healthy.
// Returns true if the server responds with HTTP 200 status, otherwise returns false.
func checkServer(client *http.Client, url string) bool {
	res, err := client.Get(url)
	if err != nil {
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	return res.StatusCode == http.StatusOK
}
