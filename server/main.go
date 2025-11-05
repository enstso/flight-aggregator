package main

import (
	"aggregator/internal/config"
	"aggregator/internal/handler"
	"aggregator/internal/health"
	"fmt"
	"net/http"
)

// withCORS is a middleware that enables CORS support for the provided HTTP handler.
// Sets headers to allow cross-origin requests and handles preflight OPTIONS requests.
func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// main initializes the server, loads configuration, defines HTTP routes, and starts listening for incoming requests.
func main() {
	config.Load()

	mux := http.NewServeMux()

	mux.HandleFunc("/health", health.HealthHandler)
	mux.HandleFunc("/flights", handler.GetFlights)
	mux.HandleFunc("/flights/id/", handler.GetFlightById)
	mux.HandleFunc("/flights/number/", handler.GetFlightByNumber)
	mux.HandleFunc("/flights/passengerName/", handler.GetFlightsByPassenger)
	mux.HandleFunc("/flights/destination", handler.GetFlightsByDestination)
	mux.HandleFunc("/flights/price/", handler.GetFlightsByPrice)
	mux.HandleFunc("/flights/sorted", handler.GetFlightsSorted)

	fmt.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", withCORS(mux)); err != nil {
		fmt.Println("Server error:", err)
	}
}
