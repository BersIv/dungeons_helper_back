package subraces

import "context"

type Subraces struct {
	Id          int64  `json:"id"`
	SubraceName string `json:"raceName"`
	IdStats     int64  `json:"idStats"`
}

type Repository interface {
	GetAllSubraces(ctx context.Context) ([]Subraces, error)
}

type Service interface {
	GetAllSubraces(c context.Context) ([]Subraces, error)
}
