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

	seg1 := domain.NewSegment(
		"AA100",
		"JFK",
		"LAX",
		now.Add(2*time.Hour),
		now.Add(7*time.Hour),
	)
	total1 := domain.NewTotal(500.00, "USD")
	f1 := domain.NewFlight(
		"1",
		"confirmed",
		"John Doe",
		[]domain.Segment{seg1},
		total1,
		"source1",
	)

	seg2 := domain.NewSegment(
		"UA200",
		"JFK",
		"SFO",
		now.Add(1*time.Hour),
		now.Add(5*time.Hour),
	)
	total2 := domain.NewTotal(300.00, "USD")
	f2 := domain.NewFlight(
		"2",
		"confirmed",
		"Jane Smith",
		[]domain.Segment{seg2},
		total2,
		"source2",
	)

	seg3 := domain.NewSegment(
		"DL300",
		"JFK",
		"ORD",
		now.Add(3*time.Hour),
		now.Add(6*time.Hour),
	)
	total3 := domain.NewTotal(400.00, "USD")
	f3 := domain.NewFlight(
		"3",
		"confirmed",
		"Bob Johnson",
		[]domain.Segment{seg3},
		total3,
		"source1",
	)

	return domain.Flights{*f1, *f2, *f3}
}

// createMockFlightsWithConnections generates mock flight data, including flights with multiple connections, for testing purposes.
func createMockFlightsWithConnections() domain.Flights {
	now := time.Now()

	seg1_1 := domain.NewSegment(
		"AA100",
		"JFK",
		"DFW",
		now.Add(2*time.Hour),
		now.Add(5*time.Hour),
	)
	seg1_2 := domain.NewSegment(
		"AA101",
		"DFW",
		"LAX",
		now.Add(6*time.Hour),
		now.Add(9*time.Hour),
	)
	total1 := domain.NewTotal(600.00, "USD")
	f1 := domain.NewFlight(
		"1",
		"confirmed",
		"John Doe",
		[]domain.Segment{seg1_1, seg1_2},
		total1,
		"source1",
	)

	seg2 := domain.NewSegment(
		"UA200",
		"JFK",
		"LAX",
		now.Add(1*time.Hour),
		now.Add(5*time.Hour),
	)
	total2 := domain.NewTotal(400.00, "USD")
	f2 := domain.NewFlight(
		"2",
		"confirmed",
		"Jane Smith",
		[]domain.Segment{seg2},
		total2,
		"source2",
	)

	return domain.Flights{*f1, *f2}
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
		assert.Equal(t, 300.00, sorted[0].Total().Amount())
		assert.Equal(t, 400.00, sorted[1].Total().Amount())
		assert.Equal(t, 500.00, sorted[2].Total().Amount())

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
		assert.Equal(t, "3", sorted[0].ID())
		assert.Equal(t, "2", sorted[1].ID())
		assert.Equal(t, "1", sorted[2].ID())

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
		assert.Equal(t, "2", sorted[0].ID())
		assert.Equal(t, "1", sorted[1].ID())

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
		assert.Equal(t, "2", sorted[0].ID())
		assert.Equal(t, "1", sorted[1].ID())
		assert.Equal(t, "3", sorted[2].ID())

		mockRepo.AssertExpectations(t)
	})

	t.Run("handles flights with no segments", func(t *testing.T) {
		now := time.Now()

		f1 := domain.NewFlight(
			"1",
			"confirmed",
			"",
			[]domain.Segment{},
			domain.NewTotal(100.00, "USD"),
			"source1",
		)

		seg2 := domain.NewSegment(
			"AA100",
			"JFK",
			"LAX",
			now,
			now.Add(5*time.Hour),
		)
		f2 := domain.NewFlight(
			"2",
			"confirmed",
			"",
			[]domain.Segment{seg2},
			domain.NewTotal(200.00, "USD"),
			"source2",
		)

		mockFlights := domain.Flights{*f1, *f2}

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
		seg := domain.NewSegment(
			"AA100",
			"JFK",
			"LAX",
			now,
			now.Add(5*time.Hour),
		)
		f := domain.NewFlight(
			"1",
			"confirmed",
			"Test",
			[]domain.Segment{seg},
			domain.NewTotal(0, "USD"),
			"source",
		)

		duration := service.TotalTravelTime(*f)

		assert.Equal(t, 5*time.Hour, duration)
	})

	t.Run("calculates travel time for multiple segments", func(t *testing.T) {
		seg1 := domain.NewSegment(
			"S1",
			"JFK",
			"ORD",
			now,
			now.Add(3*time.Hour),
		)
		seg2 := domain.NewSegment(
			"S2",
			"ORD",
			"LAX",
			now.Add(4*time.Hour),
			now.Add(7*time.Hour),
		)
		f := domain.NewFlight(
			"2",
			"confirmed",
			"Test",
			[]domain.Segment{seg1, seg2},
			domain.NewTotal(0, "USD"),
			"source",
		)

		duration := service.TotalTravelTime(*f)

		assert.Equal(t, 7*time.Hour, duration)
	})

	t.Run("returns zero for flight with no segments", func(t *testing.T) {
		f := domain.NewFlight(
			"3",
			"confirmed",
			"Test",
			[]domain.Segment{},
			domain.NewTotal(0, "USD"),
			"source",
		)

		duration := service.TotalTravelTime(*f)

		assert.Equal(t, time.Duration(0), duration)
	})
}
