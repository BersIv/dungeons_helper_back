package stats

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

func (s *service) GetStatsById(c context.Context, id int64) (*GetStatsRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.Repository.GetStatsById(ctx, id)
	if err != nil {
		return nil, err
	}

	res := &GetStatsRes{
		Strength:     r.Strength,
		Dexterity:    r.Dexterity,
		Constitution: r.Constitution,
		Intelligence: r.Intelligence,
		Wisdom:       r.Wisdom,
		Charisma:     r.Charisma,
	}

	return res, err
}
