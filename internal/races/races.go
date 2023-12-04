package races

import "context"

type Races struct {
	Id       int64  `json:"id"`
	RaceName string `json:"raceName"`
}

type Repository interface {
	GetAllRaces(ctx context.Context) ([]Races, error)
}

type Service interface {
	GetAllRaces(c context.Context) ([]Races, error)
}
