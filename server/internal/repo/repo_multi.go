package repo

import (
	"aggregator/internal/domain"
	"context"
	"errors"
)

type Multi struct {
	repos []domain.FlightsRepository
}

// NewMulti creates a new Multi instance with the provided list of FlightsRepository implementations.
func NewMulti(repos ...domain.FlightsRepository) *Multi {
	return &Multi{repos: repos}
}

// List retrieves all flights from multiple repositories and returns them as a combined collection or an error.
func (m *Multi) List(ctx context.Context) (domain.Flights, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var all domain.Flights
	for _, r := range m.repos {
		items, err := r.List(ctx)
		if err != nil {
			return nil, err
		}
		all = append(all, items...)
	}
	return all, nil
}

// FindByID searches for a flight by ID across multiple repositories and returns the flight or an error if not found.
func (m *Multi) FindByID(ctx context.Context, id string) (domain.Flight, error) {
	select {
	case <-ctx.Done():
		return domain.Flight{}, ctx.Err()
	default:
	}

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
		if f.ID() == "" {
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

// FindByNumber searches for a flight by its number across multiple repositories and returns the flight or an error if not found.
func (m *Multi) FindByNumber(ctx context.Context, id string) (domain.Flight, error) {
	select {
	case <-ctx.Done():
		return domain.Flight{}, ctx.Err()
	default:
	}

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
		if f.ID() == "" {
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

// FindByPassenger searches for flights associated with a specific passenger name across multiple repositories.
// Returns a combined collection of flights or an error if none are found.
func (m *Multi) FindByPassenger(ctx context.Context, passengerName string) (domain.Flights, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var flights domain.Flights

	for _, r := range m.repos {
		f, err := r.FindByPassenger(ctx, passengerName)
		if err != nil {
			if errors.Is(err, domain.ErrFlightNotFound) {
				continue
			}
			return domain.Flights{}, err
		}
		flights = append(flights, f...)
	}
	if len(flights) == 0 {
		return nil, domain.ErrFlightsNotFound
	}
	return flights, nil
}

// FindByDestination searches for flights across multiple repositories based on the given departure and arrival locations.
// It returns a combined collection of flights or an error if no matches are found in any repository.
func (m *Multi) FindByDestination(ctx context.Context, departure, arrival string) (domain.Flights, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var flights domain.Flights

	for _, r := range m.repos {
		f, err := r.FindByDestination(ctx, departure, arrival)
		if err != nil {
			if errors.Is(err, domain.ErrFlightNotFound) {
				continue
			}
			return domain.Flights{}, err
		}
		flights = append(flights, f...)
	}
	if len(flights) == 0 {
		return nil, domain.ErrFlightsNotFound
	}
	return flights, nil
}

// FindByPrice retrieves flights matching the specified price across multiple repositories.
// Returns a combined collection of flights or an error if none are found.
func (m *Multi) FindByPrice(ctx context.Context, price float64) (domain.Flights, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var flights domain.Flights

	for _, r := range m.repos {
		f, err := r.FindByPrice(ctx, price)
		if err != nil {
			if errors.Is(err, domain.ErrFlightNotFound) {
				continue
			}
			return domain.Flights{}, err
		}
		flights = append(flights, f...)
	}
	if len(flights) == 0 {
		return nil, domain.ErrFlightsNotFound
	}
	return flights, nil
}
