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

type RepoFlightToBook struct {
	data domain.Flights
}

// NewRepoFlightToBookFromReader creates a RepoFlightToBook by reading and decoding flight data from the provided io.Reader.
func NewRepoFlightToBookFromReader(r io.Reader) (*RepoFlightToBook, error) {
	var raw struct {
		FlightToBook []struct {
			Reference string `json:"reference"`
			Status    string `json:"status"`
			Traveler  struct {
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
			} `json:"traveler"`
			Segments []struct {
				Flight struct {
					Number string `json:"number"`
					From   string `json:"from"`
					To     string `json:"to"`
					Depart string `json:"depart"`
					Arrive string `json:"arrive"`
				} `json:"flight"`
			} `json:"segments"`
			Total struct {
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
			} `json:"total"`
		} `json:"flight_to_book"`
	}

	err := json.NewDecoder(r).Decode(&raw.FlightToBook)
	if err != nil {
		return nil, err
	}

	const layout = time.RFC3339
	out := make(domain.Flights, 0, len(raw.FlightToBook))
	for _, f := range raw.FlightToBook {
		segs := make([]domain.Segment, 0, len(f.Segments))
		for _, s := range f.Segments {
			dep, err := time.Parse(layout, s.Flight.Depart)
			if err != nil {
				return nil, fmt.Errorf("flight_to_book parse depart %w", err)
			}
			arr, err := time.Parse(layout, s.Flight.Arrive)
			if err != nil {
				return nil, fmt.Errorf("flight_to_book parse arrive %w", err)
			}

			segs = append(segs, domain.Segment{
				FlightNumber: s.Flight.Number,
				Departure:    s.Flight.From,
				Arrival:      s.Flight.To,
				DepartTime:   dep,
				ArriveTime:   arr,
			})
		}
		total := domain.Total{
			Amount:   f.Total.Amount,
			Currency: f.Total.Currency,
		}

		passenger := f.Traveler.FirstName + " " + f.Traveler.LastName

		flight := domain.NewFlight(
			f.Reference,
			f.Status,
			passenger,
			segs,
			total,
			"flight_to_book",
		)
		out = append(out, *flight)
	}
	return &RepoFlightToBook{data: out}, nil
}

// List retrieves all available flights stored in the repository. It returns a slice of flights or an error if any occurs.
func (r *RepoFlightToBook) List(ctx context.Context) (domain.Flights, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return append(domain.Flights(nil), r.data...), nil
}

// FindByNumber searches for a flight by its flight number and returns the matching flight or an empty flight if not found.
// ctx is the context for carrying deadlines and cancellation signals.
// number is the flight number to search for within the flight segments.
// Returns a domain.Flight and an error if any issue occurs during the search.
func (r *RepoFlightToBook) FindByNumber(ctx context.Context, number string) (domain.Flight, error) {
	select {
	case <-ctx.Done():
		return domain.Flight{}, ctx.Err()
	default:
	}
	for _, f := range r.data {
		for _, s := range f.Segments {
			if strings.Compare(s.FlightNumber, number) == 0 {
				return f, nil
			}
		}
	}
	return domain.Flight{}, nil
}

// FindById retrieves a flight by its ID from the repository. Returns the matching flight or an empty flight if not found.
func (r *RepoFlightToBook) FindById(ctx context.Context, id string) (domain.Flight, error) {
	select {
	case <-ctx.Done():
		return domain.Flight{}, ctx.Err()
	default:
	}
	for _, f := range r.data {
		if strings.Compare(f.ID, id) == 0 {
			return f, nil
		}
	}
	return domain.Flight{}, nil
}

// FindByPassenger retrieves flights matching the specified passenger's name from the repository. Returns flights or an error.
func (r *RepoFlightToBook) FindByPassenger(ctx context.Context, passengerName string) (domain.Flights, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	var flights domain.Flights
	for _, f := range r.data {
		if strings.Compare(f.PassengerName, passengerName) == 0 {
			flights = append(flights, f)
		}
	}
	return flights, nil
}

// FindByDestination searches for flights matching the specified departure and arrival locations and returns the results.
func (r *RepoFlightToBook) FindByDestination(ctx context.Context, departure, arrival string) (domain.Flights, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	var flights domain.Flights
	for _, f := range r.data {
		for _, seg := range f.Segments {
			if strings.Compare(seg.Departure, departure) == 0 &&
				strings.Compare(seg.Arrival, arrival) == 0 {
				flights = append(flights, f)
			}
		}
	}
	return flights, nil
}

// FindByPrice filters flights in the repository by the specified price and returns a slice of matching flights or an error.
func (r *RepoFlightToBook) FindByPrice(ctx context.Context, price float64) (domain.Flights, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	var flights []domain.Flight
	for _, f := range r.data {
		if f.Total.Amount == price {
			flights = append(flights, f)
		}
	}
	return flights, nil
}
