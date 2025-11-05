package test

import (
	"aggregator/internal/domain"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockFlightsRepository is a mock implementation using testify/mock
type MockFlightsRepository struct {
	mock.Mock
}

// List retrieves all flights from the repository. It accepts a context and returns a collection of flights and an error.
func (m *MockFlightsRepository) List(ctx context.Context) (domain.Flights, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(domain.Flights), args.Error(1)
}

// FindById retrieves a flight from the repository by its unique ID. It accepts a context and an ID, returning a flight or an error.
func (m *MockFlightsRepository) FindById(ctx context.Context, id string) (domain.Flight, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Flight), args.Error(1)
}

// FindByNumber retrieves a flight from the repository by its flight number. It accepts a context and a flight number as input.
func (m *MockFlightsRepository) FindByNumber(ctx context.Context, number string) (domain.Flight, error) {
	args := m.Called(ctx, number)
	return args.Get(0).(domain.Flight), args.Error(1)
}

// FindByPassenger retrieves flights associated with a specific passenger by their name. Returns flights or an error.
func (m *MockFlightsRepository) FindByPassenger(ctx context.Context, passengerName string) (domain.Flights, error) {
	args := m.Called(ctx, passengerName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(domain.Flights), args.Error(1)
}

// FindByDestination retrieves flights based on the given departure and arrival locations. Returns matching flights or an error.
func (m *MockFlightsRepository) FindByDestination(ctx context.Context, departure, arrival string) (domain.Flights, error) {
	args := m.Called(ctx, departure, arrival)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(domain.Flights), args.Error(1)
}

// FindByPrice retrieves flights from the repository that match the specified price. It returns a list of flights or an error.
func (m *MockFlightsRepository) FindByPrice(ctx context.Context, price float64) (domain.Flights, error) {
	args := m.Called(ctx, price)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(domain.Flights), args.Error(1)
}
