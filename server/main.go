/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"aggregator/internal/config"
	"aggregator/internal/handler"
	"aggregator/internal/health"
	"fmt"
	"net/http"
)

func main() {
	//load the config .env
	config.Load()

	http.HandleFunc("/health", health.HealthHandler)
	http.HandleFunc("/flights", handler.GetFlights)
	http.HandleFunc("/flights/id/", handler.GetFlightById)
	http.HandleFunc("/flights/number/", handler.GetFlightByNumber)
	http.HandleFunc("/flights/passengerName/", handler.GetFlightsByPassenger)
	http.HandleFunc("/flights/destination", handler.GetFlightsByDestination)
	http.HandleFunc(" /flights/price/", handler.GetFlightsByPrice)
	http.HandleFunc("/flights/sorted", handler.GetFlightsSorted)

	fmt.Println("Server running on :8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
