package class

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

func (s *service) GetAllRaces(c context.Context) ([]Class, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	classes, err := s.Repository.GetAllClasses(ctx)
	if err != nil {
		return nil, err
	}

	return classes, nil
}
