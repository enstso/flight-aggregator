package handler

import (
	"aggregator/internal/config"
	"aggregator/internal/db"
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

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(flights); err != nil {
		http.Error(w, "encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

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

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(flight); err != nil {
		http.Error(w, "encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

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

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(flight); err != nil {
		http.Error(w, "encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

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

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(flights); err != nil {
		http.Error(w, "encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
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

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(flights); err != nil {
		http.Error(w, "encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
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

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(flights); err != nil {
		http.Error(w, "encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

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

	fmt.Println("[GET] /flights/sorted type=",
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
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(flights); err != nil {
		http.Error(w, "encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

func GetMultiRepo(ctx context.Context, w http.ResponseWriter) *repo.Multi {
	b1, err := db.GetJSON(ctx, config.SERVER1_URL+"flights")
	if err != nil {
		http.Error(w, "fetch flights: "+err.Error(), http.StatusBadGateway)
		return nil
	}

	b2, err := db.GetJSON(ctx, config.SERVER2_URL+"flight_to_book")
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

type FlightDestinationRequest struct {
	Departure string `json:"departure"`
	Arrival   string `json:"arrival"`
}
