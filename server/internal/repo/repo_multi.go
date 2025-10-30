package repo

import (
	"aggregator/internal/domain"
	"context"
	"errors"
)

type Multi struct {
	repos []domain.FlightsRepository
}

func NewMulti(repos ...domain.FlightsRepository) *Multi {
	return &Multi{repos: repos}
}

// List retrieves all flights from multiple repositories and returns them as a combined collection or an error.
func (m *Multi) List(ctx context.Context) ([]domain.Flight, error) {
	var all []domain.Flight
	for _, r := range m.repos {
		items, err := r.List(ctx)
		if err != nil {
			return nil, err
		}
		all = append(all, items...)
	}
	return all, nil
}

func (m *Multi) FindByID(ctx context.Context, id string) (domain.Flight, error) {
	var lastErr error
	for _, r := range m.repos {
		f, err := r.FindById(ctx, id)
		if err != nil {
			if errors.Is(err, domain.ErrFlightNotFound) {
				lastErr = err
				continue
			}
			return domain.Flight{}, err
		}
		if f.ID == "" {
			lastErr = domain.ErrFlightNotFound
			continue
		}
		return f, nil
	}
	if lastErr == nil {
		lastErr = domain.ErrFlightNotFound
	}
	return domain.Flight{}, lastErr
}

func (m *Multi) FindByNumber(ctx context.Context, id string) (domain.Flight, error) {
	var lastErr error
	for _, r := range m.repos {
		f, err := r.FindByNumber(ctx, id)
		if err != nil {
			if errors.Is(err, domain.ErrFlightNotFound) {
				lastErr = err
				continue
			}
			return domain.Flight{}, err
		}
		if f.ID == "" {
			lastErr = domain.ErrFlightNotFound
			continue
		}
		return f, nil
	}
	if lastErr == nil {
		lastErr = domain.ErrFlightNotFound
	}
	return domain.Flight{}, lastErr
}

func (m *Multi) FindByPassenger(ctx context.Context, passengerName string) ([]domain.Flight, error) {
	var lastErr error
	for _, r := range m.repos {
		f, err := r.FindByPassenger(ctx, passengerName)
		if err != nil {
			if errors.Is(err, domain.ErrFlightNotFound) {
				lastErr = err
				continue
			}
			return domain.Flights{}, err
		}
		if len(f) == 0 {
			lastErr = domain.ErrFlightsNotFound
			continue
		}
		return f, nil
	}
	if lastErr == nil {
		lastErr = domain.ErrFlightsNotFound
	}
	return domain.Flights{}, lastErr
}

func (m *Multi) FindByDestination(ctx context.Context, departure, arrival string) ([]domain.Flight, error) {
	var lastErr error
	for _, r := range m.repos {
		f, err := r.FindByDestination(ctx, departure, arrival)
		if err != nil {
			if errors.Is(err, domain.ErrFlightNotFound) {
				lastErr = err
				continue
			}
			return domain.Flights{}, err
		}
		if len(f) == 0 {
			lastErr = domain.ErrFlightsNotFound
			continue
		}
		return f, nil
	}
	if lastErr == nil {
		lastErr = domain.ErrFlightsNotFound
	}
	return domain.Flights{}, lastErr
}

func (m *Multi) FindByPrice(ctx context.Context, price float64) ([]domain.Flight, error) {
	var lastErr error
	for _, r := range m.repos {
		f, err := r.FindByPrice(ctx, price)
		if err != nil {
			if errors.Is(err, domain.ErrFlightNotFound) {
				lastErr = err
				continue
			}
			return domain.Flights{}, err
		}
		if len(f) == 0 {
			lastErr = domain.ErrFlightsNotFound
			continue
		}
		return f, nil
	}
	if lastErr == nil {
		lastErr = domain.ErrFlightsNotFound
	}
	return domain.Flights{}, lastErr
}
