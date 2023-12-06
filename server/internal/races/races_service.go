package races

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

func (s *service) GetAllRaces(c context.Context) ([]Races, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	races, err := s.Repository.GetAllRaces(ctx)
	if err != nil {
		return nil, err
	}

	return races, nil
}
