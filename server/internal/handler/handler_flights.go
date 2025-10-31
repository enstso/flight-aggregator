package handler

import (
	"aggregator/internal/config"
	"aggregator/internal/db"
	"aggregator/internal/repo"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetFlights
func GetFlights(w http.ResponseWriter, r *http.Request) {
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
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(flight); err != nil {
		http.Error(w, "encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

func GetFlightByNumber(w http.ResponseWriter, r *http.Request) {
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
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(flight); err != nil {
		http.Error(w, "encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

func GetFlightsByPassenger(w http.ResponseWriter, r *http.Request) {
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
	var flight, err = multi.FindByPassenger(ctx, passengerName)
	if err != nil {
		http.Error(w, "flights/passengerName/:passengerName: "+err.Error(), http.StatusNotFound)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(flight); err != nil {
		http.Error(w, "encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
func GetFlightsByPrice(w http.ResponseWriter, r *http.Request) {
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
	var flight, err = multi.FindByPrice(ctx, price)
	if err != nil {
		http.Error(w, "flights/price/:price: "+err.Error(), http.StatusNotFound)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(flight); err != nil {
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
