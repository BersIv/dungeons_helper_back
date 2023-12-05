package alignment

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

func (s *service) GetAllAlignments(c context.Context) ([]Alignment, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	alignments, err := s.Repository.GetAllAlignments(ctx)
	if err != nil {
		return nil, err
	}

	return alignments, nil
}
