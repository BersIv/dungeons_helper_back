package skills

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

func (s *service) GetAllSkills(c context.Context) ([]Skills, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	skills, err := s.Repository.GetAllSkills(ctx)
	if err != nil {
		return nil, err
	}

	return skills, nil
}
