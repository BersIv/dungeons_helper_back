package subraces

import (
	"context"
	"time"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
	}
}

func (s *service) GetAllSubraces(c context.Context, idRace int64) ([]Subraces, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	races, err := s.Repository.GetAllSubraces(ctx, idRace)
	if err != nil {
		return nil, err
	}

	return races, nil
}
