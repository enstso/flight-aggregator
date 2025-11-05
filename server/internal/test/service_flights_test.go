package test

import (
	"aggregator/internal/domain"
	"aggregator/internal/repo"
	"aggregator/internal/service"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createMockFlights() domain.Flights {
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
					DepartTime:   now.Add(2 * time.Hour),
					ArriveTime:   now.Add(7 * time.Hour),
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
					DepartTime:   now.Add(1 * time.Hour),
					ArriveTime:   now.Add(5 * time.Hour),
				},
			},
			Total: domain.Total{
				Amount:   300.00,
				Currency: "USD",
			},
			Source: "source2",
		},
		{
			ID:            "3",
			Status:        "confirmed",
			PassengerName: "Bob Johnson",
			Segments: []domain.Segment{
				{
					FlightNumber: "DL300",
					Departure:    "JFK",
					Arrival:      "ORD",
					DepartTime:   now.Add(3 * time.Hour),
					ArriveTime:   now.Add(6 * time.Hour),
				},
			},
			Total: domain.Total{
				Amount:   400.00,
				Currency: "USD",
			},
			Source: "source1",
		},
	}
}

// createMockFlightsWithConnections generates mock flight data, including flights with multiple connections, for testing purposes.
func createMockFlightsWithConnections() domain.Flights {
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
					Arrival:      "DFW",
					DepartTime:   now.Add(2 * time.Hour),
					ArriveTime:   now.Add(5 * time.Hour),
				},
				{
					FlightNumber: "AA101",
					Departure:    "DFW",
					Arrival:      "LAX",
					DepartTime:   now.Add(6 * time.Hour),
					ArriveTime:   now.Add(9 * time.Hour),
				},
			},
			Total: domain.Total{
				Amount:   600.00,
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
					Arrival:      "LAX",
					DepartTime:   now.Add(1 * time.Hour),
					ArriveTime:   now.Add(5 * time.Hour),
				},
			},
			Total: domain.Total{
				Amount:   400.00,
				Currency: "USD",
			},
			Source: "source2",
		},
	}
}

// TestSortByPrice verifies the behavior of the SortByPrice function, ensuring flights are correctly sorted by price in ascending order.
func TestSortByPrice(t *testing.T) {
	println("=====================SERVICE_UNIT_TEST====================")

	ctx := context.Background()

	t.Run("sorts flights by price in ascending order", func(t *testing.T) {
		mockFlights := createMockFlights()
		mockRepo := new(MockFlightsRepository)
		mockRepo.On("List", ctx).Return(mockFlights, nil)

		multiRepo := repo.NewMulti(mockRepo)

		sorted, err := service.SortByPrice(ctx, multiRepo)

		assert.NoError(t, err)
		assert.Len(t, sorted, 3)
		assert.Equal(t, 300.00, sorted[0].Total.Amount)
		assert.Equal(t, 400.00, sorted[1].Total.Amount)
		assert.Equal(t, 500.00, sorted[2].Total.Amount)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		expectedErr := domain.ErrFlightsNotFound
		mockRepo := new(MockFlightsRepository)
		mockRepo.On("List", ctx).Return(nil, expectedErr)

		multiRepo := repo.NewMulti(mockRepo)

		_, err := service.SortByPrice(ctx, multiRepo)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("handles empty flight list", func(t *testing.T) {
		mockRepo := new(MockFlightsRepository)
		mockRepo.On("List", ctx).Return(domain.Flights{}, nil)

		multiRepo := repo.NewMulti(mockRepo)

		sorted, err := service.SortByPrice(ctx, multiRepo)

		assert.NoError(t, err)
		assert.Empty(t, sorted)

		mockRepo.AssertExpectations(t)
	})
}

// TestSortByTimeTravel validates the sorting of flights based on their total travel time in ascending order.
// It includes tests for sorting accuracy, handling of errors, empty flight lists, and flights with connections.
func TestSortByTimeTravel(t *testing.T) {
	ctx := context.Background()

	t.Run("sorts flights by total travel time", func(t *testing.T) {
		mockFlights := createMockFlights()
		mockRepo := new(MockFlightsRepository)
		mockRepo.On("List", ctx).Return(mockFlights, nil)

		multiRepo := repo.NewMulti(mockRepo)

		sorted, err := service.SortByTimeTravel(ctx, multiRepo)

		assert.NoError(t, err)
		assert.Len(t, sorted, 3)
		// Flight 3: 3h, Flight 2: 4h, Flight 1: 5h
		assert.Equal(t, "3", sorted[0].ID)
		assert.Equal(t, "2", sorted[1].ID)
		assert.Equal(t, "1", sorted[2].ID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("sorts flights with connections correctly", func(t *testing.T) {
		mockFlights := createMockFlightsWithConnections()
		mockRepo := new(MockFlightsRepository)
		mockRepo.On("List", ctx).Return(mockFlights, nil)

		multiRepo := repo.NewMulti(mockRepo)

		sorted, err := service.SortByTimeTravel(ctx, multiRepo)

		assert.NoError(t, err)
		assert.Len(t, sorted, 2)
		// Flight 2: 4h (direct), Flight 1: 7h (with connection)
		assert.Equal(t, "2", sorted[0].ID)
		assert.Equal(t, "1", sorted[1].ID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		expectedErr := domain.ErrFlightsNotFound
		mockRepo := new(MockFlightsRepository)
		mockRepo.On("List", ctx).Return(nil, expectedErr)

		multiRepo := repo.NewMulti(mockRepo)

		_, err := service.SortByTimeTravel(ctx, multiRepo)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("handles empty flight list", func(t *testing.T) {
		mockRepo := new(MockFlightsRepository)
		mockRepo.On("List", ctx).Return(domain.Flights{}, nil)

		multiRepo := repo.NewMulti(mockRepo)

		sorted, err := service.SortByTimeTravel(ctx, multiRepo)

		assert.NoError(t, err)
		assert.Empty(t, sorted)

		mockRepo.AssertExpectations(t)
	})
}

// TestSortByDepartureDate verifies the behavior of sorting flights by departure date of their segments in various scenarios.
func TestSortByDepartureDate(t *testing.T) {
	ctx := context.Background()

	t.Run("sorts flights by departure date", func(t *testing.T) {
		mockFlights := createMockFlights()
		mockRepo := new(MockFlightsRepository)
		mockRepo.On("List", ctx).Return(mockFlights, nil)

		multiRepo := repo.NewMulti(mockRepo)

		sorted, err := service.SortByDepartureDate(ctx, multiRepo)

		assert.NoError(t, err)
		assert.Len(t, sorted, 3)
		// Sorted by departure: Flight 2 (1h), Flight 1 (2h), Flight 3 (3h)
		assert.Equal(t, "2", sorted[0].ID)
		assert.Equal(t, "1", sorted[1].ID)
		assert.Equal(t, "3", sorted[2].ID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("handles flights with no segments", func(t *testing.T) {
		mockFlights := domain.Flights{
			{
				ID:       "1",
				Segments: []domain.Segment{},
				Total:    domain.Total{Amount: 100.00, Currency: "USD"},
			},
			{
				ID: "2",
				Segments: []domain.Segment{
					{
						FlightNumber: "AA100",
						Departure:    "JFK",
						Arrival:      "LAX",
						DepartTime:   time.Now(),
						ArriveTime:   time.Now().Add(5 * time.Hour),
					},
				},
				Total: domain.Total{Amount: 200.00, Currency: "USD"},
			},
		}
		mockRepo := new(MockFlightsRepository)
		mockRepo.On("List", ctx).Return(mockFlights, nil)

		multiRepo := repo.NewMulti(mockRepo)

		sorted, err := service.SortByDepartureDate(ctx, multiRepo)

		assert.NoError(t, err)
		assert.Len(t, sorted, 2)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		expectedErr := domain.ErrFlightsNotFound
		mockRepo := new(MockFlightsRepository)
		mockRepo.On("List", ctx).Return(nil, expectedErr)

		multiRepo := repo.NewMulti(mockRepo)

		_, err := service.SortByDepartureDate(ctx, multiRepo)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)

		mockRepo.AssertExpectations(t)
	})
}

// TestTotalTravelTime tests the TotalTravelTime function to ensure correct calculation of total travel time for flights.
func TestTotalTravelTime(t *testing.T) {
	now := time.Now()

	t.Run("calculates travel time for single segment", func(t *testing.T) {
		flight := domain.Flight{
			Segments: []domain.Segment{
				{
					DepartTime: now,
					ArriveTime: now.Add(5 * time.Hour),
				},
			},
		}

		duration := service.TotalTravelTime(flight)

		assert.Equal(t, 5*time.Hour, duration)
	})

	t.Run("calculates travel time for multiple segments", func(t *testing.T) {
		flight := domain.Flight{
			Segments: []domain.Segment{
				{
					DepartTime: now,
					ArriveTime: now.Add(3 * time.Hour),
				},
				{
					DepartTime: now.Add(4 * time.Hour),
					ArriveTime: now.Add(7 * time.Hour),
				},
			},
		}

		duration := service.TotalTravelTime(flight)

		assert.Equal(t, 7*time.Hour, duration)
	})

	t.Run("returns zero for flight with no segments", func(t *testing.T) {
		flight := domain.Flight{
			Segments: []domain.Segment{},
		}

		duration := service.TotalTravelTime(flight)

		assert.Equal(t, time.Duration(0), duration)
	})
}
