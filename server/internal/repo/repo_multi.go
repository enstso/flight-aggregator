package repo

import (
	"aggregator/internal/domain"
	"context"
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
