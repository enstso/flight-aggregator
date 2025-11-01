package test

import (
	"aggregator/internal/domain"
	"aggregator/internal/repo"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createTestFlights() domain.Flights {
	now := time.Now()
	return domain.Flights{
		{
			ID:            "1",
			Status:        "confirmed",
			PassengerName: "John Doe",
			Segments: []domain.Segment{
				{
					FlightNumber: "AA100",
					Departure:    "JFK",
					Arrival:      "LAX",
					DepartTime:   now,
					ArriveTime:   now.Add(5 * time.Hour),
				},
			},
			Total: domain.Total{
				Amount:   500.00,
				Currency: "USD",
			},
			Source: "source1",
		},
		{
			ID:            "2",
			Status:        "confirmed",
			PassengerName: "Jane Smith",
			Segments: []domain.Segment{
				{
					FlightNumber: "UA200",
					Departure:    "JFK",
					Arrival:      "SFO",
					DepartTime:   now,
					ArriveTime:   now.Add(4 * time.Hour),
				},
			},
			Total: domain.Total{
				Amount:   300.00,
				Currency: "USD",
			},
			Source: "source2",
		},
	}
}

func TestMulti_List(t *testing.T) {
	println("=====================REPO_UNIT_TEST====================")
	ctx := context.Background()

	t.Run("aggregates flights from multiple repositories", func(t *testing.T) {
		flights := createTestFlights()

		repo1 := new(MockFlightsRepository)
		repo1.On("List", ctx).Return(flights[:1], nil)

		repo2 := new(MockFlightsRepository)
		repo2.On("List", ctx).Return(flights[1:], nil)

		multi := repo.NewMulti(repo1, repo2)

		result, err := multi.List(ctx)

		assert.NoError(t, err)
		assert.Len(t, result, 2)

		repo1.AssertExpectations(t)
		repo2.AssertExpectations(t)
	})

	t.Run("returns error when any repository fails", func(t *testing.T) {
		expectedErr := errors.New("repository error")
		flights := createTestFlights()

		repo1 := new(MockFlightsRepository)
		repo1.On("List", ctx).Return(flights[:1], nil)

		repo2 := new(MockFlightsRepository)
		repo2.On("List", ctx).Return(nil, expectedErr)

		multi := repo.NewMulti(repo1, repo2)

		_, err := multi.List(ctx)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)

		repo1.AssertExpectations(t)
		repo2.AssertExpectations(t)
	})

	t.Run("handles empty repositories", func(t *testing.T) {
		repo1 := new(MockFlightsRepository)
		repo1.On("List", ctx).Return(domain.Flights{}, nil)

		repo2 := new(MockFlightsRepository)
		repo2.On("List", ctx).Return(domain.Flights{}, nil)

		multi := repo.NewMulti(repo1, repo2)

		result, err := multi.List(ctx)

		assert.NoError(t, err)
		assert.Empty(t, result)

		repo1.AssertExpectations(t)
		repo2.AssertExpectations(t)
	})
}

func TestMulti_FindByID(t *testing.T) {
	ctx := context.Background()

	t.Run("finds flight by ID from first repository", func(t *testing.T) {
		flights := createTestFlights()

		repo1 := new(MockFlightsRepository)
		repo1.On("FindById", ctx, "1").Return(flights[0], nil)

		repo2 := new(MockFlightsRepository)

		multi := repo.NewMulti(repo1, repo2)

		result, err := multi.FindByID(ctx, "1")

		assert.NoError(t, err)
		assert.Equal(t, "1", result.ID)

		repo1.AssertExpectations(t)
	})

	t.Run("finds flight by ID from second repository", func(t *testing.T) {
		flights := createTestFlights()

		repo1 := new(MockFlightsRepository)
		repo1.On("FindById", ctx, "2").Return(domain.Flight{}, domain.ErrFlightNotFound)

		repo2 := new(MockFlightsRepository)
		repo2.On("FindById", ctx, "2").Return(flights[1], nil)

		multi := repo.NewMulti(repo1, repo2)

		result, err := multi.FindByID(ctx, "2")

		assert.NoError(t, err)
		assert.Equal(t, "2", result.ID)

		repo1.AssertExpectations(t)
		repo2.AssertExpectations(t)
	})

	t.Run("returns error when flight not found", func(t *testing.T) {
		repo1 := new(MockFlightsRepository)
		repo1.On("FindById", ctx, "999").Return(domain.Flight{}, domain.ErrFlightNotFound)

		repo2 := new(MockFlightsRepository)
		repo2.On("FindById", ctx, "999").Return(domain.Flight{}, domain.ErrFlightNotFound)

		multi := repo.NewMulti(repo1, repo2)

		_, err := multi.FindByID(ctx, "999")

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrFlightNotFound)

		repo1.AssertExpectations(t)
		repo2.AssertExpectations(t)
	})

	t.Run("handles repository error", func(t *testing.T) {
		expectedErr := errors.New("database error")

		repo1 := new(MockFlightsRepository)
		repo1.On("FindById", ctx, "1").Return(domain.Flight{}, expectedErr)

		multi := repo.NewMulti(repo1)

		_, err := multi.FindByID(ctx, "1")

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)

		repo1.AssertExpectations(t)
	})
}

func TestMulti_FindByNumber(t *testing.T) {
	ctx := context.Background()

	t.Run("finds flight by number", func(t *testing.T) {
		flights := createTestFlights()

		repo1 := new(MockFlightsRepository)
		repo1.On("FindByNumber", ctx, "AA100").Return(flights[0], nil)

		multi := repo.NewMulti(repo1)

		result, err := multi.FindByNumber(ctx, "AA100")

		assert.NoError(t, err)
		assert.Equal(t, "AA100", result.Segments[0].FlightNumber)

		repo1.AssertExpectations(t)
	})

	t.Run("returns error when flight not found", func(t *testing.T) {
		repo1 := new(MockFlightsRepository)
		repo1.On("FindByNumber", ctx, "XX999").Return(domain.Flight{}, domain.ErrFlightNotFound)

		multi := repo.NewMulti(repo1)

		_, err := multi.FindByNumber(ctx, "XX999")

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrFlightNotFound)

		repo1.AssertExpectations(t)
	})
}

func TestMulti_FindByPassenger(t *testing.T) {
	ctx := context.Background()

	t.Run("finds flights by passenger name", func(t *testing.T) {
		flights := createTestFlights()

		repo1 := new(MockFlightsRepository)
		repo1.On("FindByPassenger", ctx, "John Doe").Return(flights[:1], nil)

		repo2 := new(MockFlightsRepository)
		repo2.On("FindByPassenger", ctx, "John Doe").Return(nil, domain.ErrFlightNotFound)

		multi := repo.NewMulti(repo1, repo2)

		result, err := multi.FindByPassenger(ctx, "John Doe")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "John Doe", result[0].PassengerName)

		repo1.AssertExpectations(t)
		repo2.AssertExpectations(t)
	})

	t.Run("returns error when no flights found", func(t *testing.T) {
		repo1 := new(MockFlightsRepository)
		repo1.On("FindByPassenger", ctx, "Unknown Person").Return(nil, domain.ErrFlightNotFound)

		multi := repo.NewMulti(repo1)

		_, err := multi.FindByPassenger(ctx, "Unknown Person")

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrFlightsNotFound)

		repo1.AssertExpectations(t)
	})

	t.Run("aggregates flights from multiple repositories", func(t *testing.T) {
		flights := createTestFlights()

		repo1 := new(MockFlightsRepository)
		repo1.On("FindByPassenger", ctx, "Test User").Return(flights[:1], nil)

		repo2 := new(MockFlightsRepository)
		repo2.On("FindByPassenger", ctx, "Test User").Return(flights[1:], nil)

		multi := repo.NewMulti(repo1, repo2)

		result, err := multi.FindByPassenger(ctx, "Test User")

		assert.NoError(t, err)
		assert.Len(t, result, 2)

		repo1.AssertExpectations(t)
		repo2.AssertExpectations(t)
	})
}

func TestMulti_FindByDestination(t *testing.T) {
	ctx := context.Background()

	t.Run("finds flights by destination", func(t *testing.T) {
		flights := createTestFlights()

		repo1 := new(MockFlightsRepository)
		repo1.On("FindByDestination", ctx, "JFK", "LAX").Return(flights[:1], nil)

		repo2 := new(MockFlightsRepository)
		repo2.On("FindByDestination", ctx, "JFK", "LAX").Return(nil, domain.ErrFlightNotFound)

		multi := repo.NewMulti(repo1, repo2)

		result, err := multi.FindByDestination(ctx, "JFK", "LAX")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "JFK", result[0].Segments[0].Departure)
		assert.Equal(t, "LAX", result[0].Segments[0].Arrival)

		repo1.AssertExpectations(t)
		repo2.AssertExpectations(t)
	})

	t.Run("returns error when no flights found", func(t *testing.T) {
		repo1 := new(MockFlightsRepository)
		repo1.On("FindByDestination", ctx, "JFK", "XXX").Return(nil, domain.ErrFlightNotFound)

		multi := repo.NewMulti(repo1)

		_, err := multi.FindByDestination(ctx, "JFK", "XXX")

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrFlightsNotFound)

		repo1.AssertExpectations(t)
	})
}

func TestMulti_FindByPrice(t *testing.T) {
	ctx := context.Background()

	t.Run("finds flights by price", func(t *testing.T) {
		flights := createTestFlights()

		repo1 := new(MockFlightsRepository)
		repo1.On("FindByPrice", ctx, 500.00).Return(flights[:1], nil)

		repo2 := new(MockFlightsRepository)
		repo2.On("FindByPrice", ctx, 500.00).Return(nil, domain.ErrFlightNotFound)

		multi := repo.NewMulti(repo1, repo2)

		result, err := multi.FindByPrice(ctx, 500.00)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, 500.00, result[0].Total.Amount)

		repo1.AssertExpectations(t)
		repo2.AssertExpectations(t)
	})

	t.Run("returns error when no flights found", func(t *testing.T) {
		repo1 := new(MockFlightsRepository)
		repo1.On("FindByPrice", ctx, 999.99).Return(nil, domain.ErrFlightNotFound)

		multi := repo.NewMulti(repo1)

		_, err := multi.FindByPrice(ctx, 999.99)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrFlightsNotFound)

		repo1.AssertExpectations(t)
	})

	t.Run("aggregates flights with same price from multiple repositories", func(t *testing.T) {
		flights := createTestFlights()

		repo1 := new(MockFlightsRepository)
		repo1.On("FindByPrice", ctx, 500.00).Return(flights[:1], nil)

		repo2 := new(MockFlightsRepository)
		repo2.On("FindByPrice", ctx, 500.00).Return(flights[1:], nil)

		multi := repo.NewMulti(repo1, repo2)

		result, err := multi.FindByPrice(ctx, 500.00)

		assert.NoError(t, err)
		assert.Len(t, result, 2)

		repo1.AssertExpectations(t)
		repo2.AssertExpectations(t)
	})
}
