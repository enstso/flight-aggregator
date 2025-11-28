package handler

import (
	"aggregator/internal/api"
	"aggregator/internal/config"
	"aggregator/internal/domain"
	"aggregator/internal/repo"
	"aggregator/internal/service"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var errNotAllowed = errors.New("method not allowed")

// GetFlights is an HTTP handler that retrieves and returns a list of flights in JSON format for GET requests.
// Responds with an error if the method is not GET or if any issues occur during the processing.
func GetFlights(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, errNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("[GET] /flights", time.Now().Format("2006-01-02 15:04:05"))

	multi := GetMultiRepo(ctx, w)

	flights, err := multi.List(ctx)
	if err != nil {
		http.Error(w, "list flights: "+err.Error(), http.StatusInternalServerError)
		return
	}
	snapshot := flights.ToSnapshot()

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(snapshot); err != nil {
		http.Error(w, "encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

// GetFlightById handles HTTP GET requests to retrieve a flight by its unique ID from the endpoints repository system.
// It validates the HTTP method, processes the request context, and fetches flight data for a given ID.
// Returns the flight details in JSON format or an appropriate HTTP error status in case of failure.
func GetFlightById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, errNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	var ctx = r.Context()
	w.Header().Set("Content-Type", "application/json")

	var parts = strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "ID is not provided", http.StatusBadRequest)
		return
	}

	var id = parts[3]
	fmt.Println("[GET] /flights/id/", id, time.Now().Format("2006-01-02 15:04:05"))

	multi := GetMultiRepo(ctx, w)

	var flight, err = multi.FindByID(ctx, id)
	if err != nil {
		http.Error(w, "flights/id/:id: "+err.Error(), http.StatusNotFound)
		return
	}

	snapshot := flight.Snapshot()

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(snapshot); err != nil {
		http.Error(w, "encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

// GetFlightByNumber handles HTTP GET requests to retrieve flight details by its number from multiple repositories.
// Returns flight details in JSON format or an appropriate HTTP error response if the flight is not found.
// Expects the flight number as part of the URL path in the format "/flights/number/{flightNumber}".
func GetFlightByNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, errNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	var ctx = r.Context()
	w.Header().Set("Content-Type", "application/json")

	var parts = strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Number flight is not provided", http.StatusBadRequest)
		return
	}

	var number = parts[3]
	fmt.Println("[GET] /flights/number/", number, time.Now().Format("2006-01-02 15:04:05"))

	multi := GetMultiRepo(ctx, w)
	var flight, err = multi.FindByNumber(ctx, number)
	if err != nil {
		http.Error(w, "flights/number/:number: "+err.Error(), http.StatusNotFound)
		return
	}

	snapshot := flight.Snapshot()

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(snapshot); err != nil {
		http.Error(w, "encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

// GetFlightsByPassenger handles retrieving flights based on a given passenger's name from multiple repositories.
// It accepts only GET requests and expects the passenger's name in the URL path as the fourth segment.
// Returns a JSON response with a list of flights or an error message in case of failure.
func GetFlightsByPassenger(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, errNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	var ctx = r.Context()
	w.Header().Set("Content-Type", "application/json")

	var parts = strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "PassengerName is not provided", http.StatusBadRequest)
		return
	}

	var passengerName = parts[3]
	fmt.Println("[GET] /flights/passengerName/", passengerName, time.Now().Format("2006-01-02 15:04:05"))

	multi := GetMultiRepo(ctx, w)
	var flights, err = multi.FindByPassenger(ctx, passengerName)
	if err != nil {
		http.Error(w, "flights/passengerName/:passengerName: "+err.Error(), http.StatusNotFound)
		return
	}

	snapshot := flights.ToSnapshot()

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(snapshot); err != nil {
		http.Error(w, "encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

// GetFlightsByDestination handles GET requests to retrieve flights based on departure and arrival destinations.
// It expects a JSON payload containing "departure" and "arrival" fields and returns matching flights in JSON format.
// If the method is not GET, it responds with a "method not allowed" error.
// The function limits the request body size to 1MB and ensures only valid JSON is processed.
// It uses a multi-repository to search for flights and returns an error if no matches are found or on processing failures.
func GetFlightsByDestination(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, errNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	var ctx = r.Context()
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("[GET] /flights/destination", time.Now().Format("2006-01-02 15:04:05"))
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req FlightDestinationRequest
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	// ensure there's no extra JSON after the object
	if dec.More() {
		http.Error(w, "invalid JSON: multiple JSON values", http.StatusBadRequest)
		return
	}

	multi := GetMultiRepo(ctx, w)

	var flights, err = multi.FindByDestination(ctx, req.Departure, req.Arrival)
	if err != nil {
		http.Error(w, "flights/destination "+err.Error(), http.StatusNotFound)
		return
	}

	snapshot := flights.ToSnapshot()

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(snapshot); err != nil {
		http.Error(w, "encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

// GetFlightsByPrice handles HTTP GET requests to fetch flights filtered by a specified price.
// It extracts the price from the URL path, queries multiple repositories, and returns matching flights in JSON format.
// Responds with appropriate HTTP status codes for errors like bad requests, method not allowed, or data not found.
func GetFlightsByPrice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, errNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	var ctx = r.Context()
	w.Header().Set("Content-Type", "application/json")

	var parts = strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "PassengerName is not provided", http.StatusBadRequest)
		return
	}

	var priceStr = parts[3]
	fmt.Println("[GET] /flights/price/", priceStr, time.Now().Format("2006-01-02 15:04:05"))

	var price, _ = strconv.ParseFloat(priceStr, 64)

	multi := GetMultiRepo(ctx, w)

	var flights, err = multi.FindByPrice(ctx, price)
	if err != nil {
		http.Error(w, "flights/price/:price: "+err.Error(), http.StatusNotFound)
		return
	}

	snapshot := flights.ToSnapshot()

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(snapshot); err != nil {
		http.Error(w, "encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

// GetFlightsSorted handles HTTP GET requests to return a list of flights sorted by a specified type (e.g., price, time).
// It validates the method, parses the query parameter for sorting type, fetches flight data, and sorts accordingly.
// Supported sorting types include "price", "time", and "departure". Responds with JSON on success or an error message on failure.
func GetFlightsSorted(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, errNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	sortType := strings.ToLower(r.URL.Query().Get("type"))

	if sortType == "" {
		http.Error(w, "missing query param: type", http.StatusBadRequest)
		return
	}

	fmt.Println("[GET] /flights/sorted?type=",
		sortType, time.Now().Format("2006-01-02 15:04:05"))

	multi := GetMultiRepo(ctx, w)
	if multi == nil {
		return
	}
	var (
		flights domain.Flights
		err     error
	)

	switch sortType {
	case "price":
		flights, err = service.SortByPrice(ctx, multi)
	case "time", "timetravel", "duration":
		flights, err = service.SortByTimeTravel(ctx, multi)
	case "departure", "depart", "departure_date":
		flights, err = service.SortByDepartureDate(ctx, multi)
	default:
		http.Error(w, "invalid sort type: "+sortType, http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "sort flights: "+err.Error(), http.StatusInternalServerError)
		return
	}

	snapshot := flights.ToSnapshot()

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(snapshot); err != nil {
		http.Error(w, "encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

func GetMultiRepo(ctx context.Context, w http.ResponseWriter) *repo.Multi {
	b1, err := api.GetDataFromApi(ctx, config.SERVER1_URL+"flights")
	if err != nil {
		http.Error(w, "fetch flights: "+err.Error(), http.StatusBadGateway)
		return nil
	}

	b2, err := api.GetDataFromApi(ctx, config.SERVER2_URL+"flight_to_book")
	if err != nil {
		http.Error(w, "fetch flight_to_book: "+err.Error(), http.StatusBadGateway)
		return nil
	}

	rA, err := repo.NewRepoFlightsFromReader(bytes.NewReader(b1))
	if err != nil {
		http.Error(w, "decode flights: "+err.Error(), http.StatusInternalServerError)
		return nil
	}
	rB, err := repo.NewRepoFlightToBookFromReader(bytes.NewReader(b2))
	if err != nil {
		http.Error(w, "decode flight_to_book: "+err.Error(), http.StatusInternalServerError)
		return nil
	}

	return repo.NewMulti(rA, rB)
}

// FlightDestinationRequest represents a request for searching flights based on departure and arrival locations.
type FlightDestinationRequest struct {
	Departure string `json:"departure"`
	Arrival   string `json:"arrival"`
}
