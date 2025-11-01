package service

import (
	"aggregator/internal/domain"
	"aggregator/internal/repo"
	"context"
	"sort"
	"time"
)

func SortByPrice(ctx context.Context, r *repo.Multi) (domain.Flights, error) {
	var sortFlights, err = r.List(ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(sortFlights, func(i, j int) bool {
		return sortFlights[i].Total.Amount <
			sortFlights[j].Total.Amount
	})
	return sortFlights, nil
}

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

func SortByDepartureDate(ctx context.Context, r *repo.Multi) (domain.Flights, error) {
	flights, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(flights, func(i, j int) bool {
		if len(flights[i].Segments) == 0 || len(flights[j].Segments) == 0 {
			return false
		}
		return flights[i].Segments[0].DepartTime.Before(flights[j].Segments[0].DepartTime)
	})

	return flights, nil
}

func TotalTravelTime(f domain.Flight) time.Duration {
	if len(f.Segments) == 0 {
		return 0
	}
	firstDepart := f.Segments[0].DepartTime
	lastArrive := f.Segments[len(f.Segments)-1].ArriveTime
	return lastArrive.Sub(firstDepart)
}
