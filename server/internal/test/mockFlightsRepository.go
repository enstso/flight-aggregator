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

func (m *MockFlightsRepository) List(ctx context.Context) (domain.Flights, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(domain.Flights), args.Error(1)
}

func (m *MockFlightsRepository) FindById(ctx context.Context, id string) (domain.Flight, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Flight), args.Error(1)
}

func (m *MockFlightsRepository) FindByNumber(ctx context.Context, number string) (domain.Flight, error) {
	args := m.Called(ctx, number)
	return args.Get(0).(domain.Flight), args.Error(1)
}

func (m *MockFlightsRepository) FindByPassenger(ctx context.Context, passengerName string) (domain.Flights, error) {
	args := m.Called(ctx, passengerName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(domain.Flights), args.Error(1)
}

func (m *MockFlightsRepository) FindByDestination(ctx context.Context, departure, arrival string) (domain.Flights, error) {
	args := m.Called(ctx, departure, arrival)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(domain.Flights), args.Error(1)
}

func (m *MockFlightsRepository) FindByPrice(ctx context.Context, price float64) (domain.Flights, error) {
	args := m.Called(ctx, price)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(domain.Flights), args.Error(1)
}
