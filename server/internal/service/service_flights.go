package service

import (
	"aggregator/internal/domain"
	"aggregator/internal/repo"
	"context"
	"sort"
	"time"
)

// SortByPrice retrieves a list of flights from repositories and sorts them in ascending order by their total price.
func SortByPrice(ctx context.Context, r *repo.Multi) (domain.Flights, error) {
	sortFlights, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(sortFlights, func(i, j int) bool {
		return sortFlights[i].Total().Amount() <
			sortFlights[j].Total().Amount()
	})
	return sortFlights, nil
}

// SortByTimeTravel retrieves and sorts flights by their total travel time in ascending order.
// Returns the sorted flights or an error if it fails to retrieve or sort the flights.
func SortByTimeTravel(ctx context.Context, r *repo.Multi) (domain.Flights, error) {
	flights, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(flights, func(i, j int) bool {
		durationI := TotalTravelTime(flights[i])
		durationJ := TotalTravelTime(flights[j])
		return durationI < durationJ
	})
	return flights, nil
}

// SortByDepartureDate retrieves flights and sorts them by the earliest departure date of their segments.
// It returns the sorted list of flights or an error.
func SortByDepartureDate(ctx context.Context, r *repo.Multi) (domain.Flights, error) {
	flights, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(flights, func(i, j int) bool {
		segmentsI := flights[i].Segments()
		segmentsJ := flights[j].Segments()

		if len(segmentsI) == 0 || len(segmentsJ) == 0 {
			return false
		}
		return segmentsI[0].DepartTime().Before(segmentsJ[0].DepartTime())
	})

	return flights, nil
}

// TotalTravelTime calculates the total travel time of a flight by measuring the time difference between the first departure and last arrival.
// Returns zero if the flight has no segments.
func TotalTravelTime(f domain.Flight) time.Duration {
	segs := f.Segments()
	if len(segs) == 0 {
		return 0
	}
	firstDepart := segs[0].DepartTime()
	lastArrive := segs[len(segs)-1].ArriveTime()
	return lastArrive.Sub(firstDepart)
}
