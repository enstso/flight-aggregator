package repo

import (
	"aggregator/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type RepoFlights struct {
	data domain.Flights
}

// NewRepoFlightsFromReader parses flight data from an io.Reader and returns a RepoFlights instance or an error.
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

	err := json.NewDecoder(r).Decode(&raw.Flights)
	if err != nil {
		return nil, err
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

		seg := domain.NewSegment(
			f.FlightNumber,
			f.DepartureAirport,
			f.ArrivalAirport,
			dep,
			arr,
		)

		total := domain.NewTotal(
			f.Price,
			f.Currency,
		)

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

// List retrieves all flights currently stored in the repository as a domain.Flights collection.
func (r *RepoFlights) List(ctx context.Context) (domain.Flights, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return append(domain.Flights(nil), r.data...), nil
}

// FindById retrieves a flight by its unique identifier from the repository. Returns the flight or an empty object if not found.
func (r *RepoFlights) FindById(ctx context.Context, id string) (domain.Flight, error) {
	select {
	case <-ctx.Done():
		return domain.Flight{}, ctx.Err()
	default:
	}
	for _, f := range r.data {
		if strings.Compare(f.ID(), id) == 0 {
			return f, nil
		}
	}
	return domain.Flight{}, nil
}

// FindByNumber retrieves a flight from the repository by the provided flight number. Returns the flight or an empty object if not found.
func (r *RepoFlights) FindByNumber(ctx context.Context, number string) (domain.Flight, error) {
	select {
	case <-ctx.Done():
		return domain.Flight{}, ctx.Err()
	default:
	}
	for _, f := range r.data {
		for _, s := range f.Segments() {
			if strings.Compare(s.FlightNumber(), number) == 0 {
				return f, nil
			}
		}
	}
	return domain.Flight{}, nil
}

// FindByPassenger retrieves all flights that match the specified passenger name from the repository. Returns the flights or an empty collection.
func (r *RepoFlights) FindByPassenger(ctx context.Context, passengerName string) (domain.Flights, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	var flights []domain.Flight
	for _, f := range r.data {
		if strings.Compare(f.PassengerName(), passengerName) == 0 {
			flights = append(flights, f)
		}
	}
	return flights, nil
}

// FindByDestination retrieves flights that match the specified departure and arrival locations from the repository.
func (r *RepoFlights) FindByDestination(ctx context.Context, departure, arrival string) (domain.Flights, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	var flights []domain.Flight
	for _, f := range r.data {
		for _, seg := range f.Segments() {
			if strings.Compare(seg.Departure(), departure) == 0 &&
				strings.Compare(seg.Arrival(), arrival) == 0 {
				flights = append(flights, f)
			}
		}
	}
	return flights, nil
}

// FindByPrice retrieves all flights whose total price matches the specified amount. Returns the flights or an empty collection.
func (r *RepoFlights) FindByPrice(ctx context.Context, price float64) (domain.Flights, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	var flights []domain.Flight
	for _, f := range r.data {
		if f.Total().Amount() == price {
			flights = append(flights, f)
		}
	}
	return flights, nil
}
