package repo

import (
	"aggregator/internal/domain"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type RepoFlights struct {
	data domain.Flights
}

func NewRepoFlightsFromReader(r io.Reader) (*RepoFlights, error) {
	var raw struct {
		Flights []struct {
			BookingID        string  `json:"bookingId"`
			Status           string  `json:"status"`
			PassengerName    string  `json:"passengerName"`
			FlightNumber     string  `json:"flightNumber"`
			DepartureAirport string  `json:"departureAirport"`
			ArrivalAirport   string  `json:"arrivalAirport"`
			DepartureTime    string  `json:"departureTime"`
			ArrivalTime      string  `json:"arrivalTime"`
			Price            float64 `json:"price"`
			Currency         string  `json:"currency"`
		} `json:"flights"`
	}
	if err := json.NewDecoder(r).Decode(&raw); err != nil {
		return nil, fmt.Errorf("repoA decode: %w", err)
	}
	const layout = time.RFC3339
	out := make(domain.Flights, 0, len(raw.Flights))

	for _, f := range raw.Flights {
		dep, err := time.Parse(layout, f.DepartureTime)
		if err != nil {
			return nil, fmt.Errorf("repo Flights parse depart %w", err)
		}
		arr, err := time.Parse(layout, f.ArrivalTime)
		if err != nil {
			return nil, fmt.Errorf("repo Flights parse arrival %w", err)
		}

		seg := domain.Segment{
			FlightNumber: f.FlightNumber,
			Departure:    f.DepartureAirport,
			Arrival:      f.ArrivalAirport,
			DepartTime:   dep,
			ArriveTime:   arr,
		}

		total := domain.Total{
			Amount:   f.Price,
			Currency: f.Currency,
		}

		flight := domain.NewFlight(
			f.BookingID,
			f.Status,
			f.PassengerName,
			[]domain.Segment{seg},
			total,
			"flights",
		)
		out = append(out, *flight)
	}
	return &RepoFlights{out}, nil
}
