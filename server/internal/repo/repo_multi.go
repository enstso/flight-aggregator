package repo

import "aggregator/internal/domain"

type Multi struct {
	repos []domain.FlightsRepository
}

func NewMulti(repos ...domain.FlightsRepository) *Multi {
	return &Multi{repos: repos}
}
