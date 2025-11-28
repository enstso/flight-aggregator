package domain

import (
	"context"
	"errors"
	"time"
)

type Total struct {
	amount   float64
	currency string
}

type Segment struct {
	flightNumber string
	departure    string
	arrival      string
	departTime   time.Time
	arriveTime   time.Time
}

type Flight struct {
	id            string
	status        string
	passengerName string
	segments      []Segment
	total         Total
	source        string
}

type Flights []Flight

// Getters to access to private properties
func (t Total) Amount() float64  { return t.amount }
func (t Total) Currency() string { return t.currency }

func (s Segment) FlightNumber() string  { return s.flightNumber }
func (s Segment) Departure() string     { return s.departure }
func (s Segment) Arrival() string       { return s.arrival }
func (s Segment) DepartTime() time.Time { return s.departTime }
func (s Segment) ArriveTime() time.Time { return s.arriveTime }

func (f Flight) ID() string            { return f.id }
func (f Flight) Status() string        { return f.status }
func (f Flight) PassengerName() string { return f.passengerName }
func (f Flight) Segments() []Segment   { return append([]Segment(nil), f.segments...) }
func (f Flight) Total() Total          { return f.total }
func (f Flight) Source() string        { return f.source }

// NewFlight creates and returns a new Flight instance with the specified ID, status, passenger name, segments, total, and source.
func NewFlight(id, status, passengerName string, segments []Segment, total Total, source string) *Flight {
	return &Flight{
		id:            id,
		status:        status,
		passengerName: passengerName,
		segments:      append([]Segment(nil), segments...),
		total:         total,
		source:        source,
	}
}

func (t TotalSnapshot) ToDomain() Total {
	return Total{t.Amount, t.Currency}
}

func (s SegmentSnapshot) ToDomain() Segment {
	return Segment{
		s.FlightNumber,
		s.Departure,
		s.Arrival,
		s.DepartTime,
		s.ArriveTime,
	}
}

func (f FlightSnapshot) ToDomain() *Flight {
	segs := make([]Segment, len(f.Segments))
	for i, s := range f.Segments {
		segs[i] = s.ToDomain()
	}
	return NewFlight(f.ID, f.Status, f.PassengerName, segs, f.Total.ToDomain(), f.Source)
}

func (fs FlightsSnapshot) ToDomain() Flights {
	out := make(Flights, len(fs))
	for i, s := range fs {
		out[i] = *s.ToDomain()
	}
	return out
}

// FlightsRepository defines methods to interact with flight data and perform various queries.
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

// NewTotal (unitTest) creates a new Total instance with the specified amount and currency.
func NewTotal(amount float64, currency string) Total {
	return Total{amount: amount, currency: currency}
}

// NewSegment (unitTest) creates and returns a new Segment with the provided flight number, departure, arrival, and timing details.
func NewSegment(flightNumber, departure, arrival string, depart, arrive time.Time) Segment {
	return Segment{
		flightNumber: flightNumber,
		departure:    departure,
		arrival:      arrival,
		departTime:   depart,
		arriveTime:   arrive,
	}
}

var ErrFlightNotFound = errors.New("flight not found")
var ErrFlightsNotFound = errors.New("flights not found")
