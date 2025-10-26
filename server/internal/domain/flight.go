package domain

import (
	"context"
	"time"
)

type Total struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type Segment struct {
	FlightNumber string    `json:"flightNumber"`
	Departure    string    `json:"from"`
	Arrival      string    `json:"to"`
	DepartTime   time.Time `json:"depart"`
	ArriveTime   time.Time `json:"arrive"`
}

type Flight struct {
	ID            string    `json:"id"`
	Status        string    `json:"status"`
	PassengerName string    `json:"passengerName"`
	Segments      []Segment `json:"segments"`
	Total         Total     `json:"total"`
	Source        string    `json:"source"`
}

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

type FlightsRepository interface {
	List(ctx context.Context) (Flights, error)

	FindByNumber(ctx context.Context, number string) (Flight, error)

	FindByPassenger(ctx context.Context, firstName, lastName string) (Flights, error)

	FindByDestination(ctx context.Context, departure, arrival string) (Flights, error)

	FindByPrice(ctx context.Context, price float64) (Flights, error)
}
