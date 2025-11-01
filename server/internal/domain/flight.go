package domain

import (
	"context"
	"errors"
	"time"
)

// Total represents an amount and its associated currency.
type Total struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// Segment represents a flight segment with details of departure, arrival, and timing information.
type Segment struct {
	FlightNumber string    `json:"flightNumber"`
	Departure    string    `json:"from"`
	Arrival      string    `json:"to"`
	DepartTime   time.Time `json:"depart"`
	ArriveTime   time.Time `json:"arrive"`
}

// Flight represents a flight containing its ID, status, passenger information, segments, total, and source details.
type Flight struct {
	ID            string    `json:"id"`
	Status        string    `json:"status"`
	PassengerName string    `json:"passengerName"`
	Segments      []Segment `json:"segments"`
	Total         Total     `json:"total"`
	Source        string    `json:"source"`
}

// NewFlight creates and returns a new Flight instance with the specified ID, status, passenger name, segments, total, and source.
func NewFlight(id, status, passengerName string,
	segments []Segment, total Total, source string) *Flight {
	return &Flight{
		id,
		status,
		passengerName,
		segments,
		total,
		source,
	}
}

type Flights []Flight

// FlightsRepository defines methods to interact with flight data and perform various queries.
// List retrieves all flights.
// FindById retrieves a flight by its unique ID from the repository.
// FindByNumber retrieves a flight by its flight number.
// FindByPassenger retrieves flights based on a passenger's first and last name.
// FindByDestination retrieves flights from a departure location to an arrival location.
// FindByPrice retrieves flights with a specified price.
type FlightsRepository interface {

	// List retrieves all available flights from the repository and returns them as a collection.
	List(ctx context.Context) (Flights, error)

	// FindById retrieves a flight by its unique ID from the repository.
	FindById(ctx context.Context, id string) (Flight, error)

	// FindByNumber retrieves a specific flight from the repository based on the provided flight number.
	FindByNumber(ctx context.Context, number string) (Flight, error)

	// FindByPassenger retrieves flights associated with a specific passenger's first and last name (passengerName) from the repository.
	FindByPassenger(ctx context.Context, passengerName string) (Flights, error)

	// FindByDestination retrieves flights that match the specified departure and arrival locations from the repository.
	FindByDestination(ctx context.Context, departure, arrival string) (Flights, error)

	// FindByPrice retrieves flights from the repository that match the specified price.
	FindByPrice(ctx context.Context, price float64) (Flights, error)
}

var ErrFlightNotFound = errors.New("flight not found")
var ErrFlightsNotFound = errors.New("flights not found")
