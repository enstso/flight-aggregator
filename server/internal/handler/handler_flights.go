package handler

import (
	"aggregator/internal/config"
	"aggregator/internal/db"
	"aggregator/internal/repo"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GetFlights est un handler HTTP qui agrège toutes les sources de vols.
func GetFlights(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("[GET] /flights", time.Now().Format("2006-01-02 15:04:05"))
	b1, err := db.GetJSON(ctx, config.SERVER1_URL+"flights")
	if err != nil {
		http.Error(w, "fetch flights: "+err.Error(), http.StatusBadGateway)
		return
	}

	b2, err := db.GetJSON(ctx, config.SERVER2_URL+"flight_to_book")
	if err != nil {
		http.Error(w, "fetch flight_to_book: "+err.Error(), http.StatusBadGateway)
		return
	}

	rA, err := repo.NewRepoFlightsFromReader(bytes.NewReader(b1))
	if err != nil {
		http.Error(w, "decode flights: "+err.Error(), http.StatusInternalServerError)
		return
	}
	rB, err := repo.NewRepoFlightToBookFromReader(bytes.NewReader(b2))
	if err != nil {
		http.Error(w, "decode flight_to_book: "+err.Error(), http.StatusInternalServerError)
		return
	}

	multi := repo.NewMulti(rA, rB)

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
