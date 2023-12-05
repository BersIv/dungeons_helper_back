package stats

import "context"

type Stats struct {
	Id           int64 `json:"id"`
	Strength     int64 `json:"strength"`
	Dexterity    int64 `json:"dexterity"`
	Constitution int64 `json:"constitution"`
	Intelligence int64 `json:"intelligence"`
	Wisdom       int64 `json:"wisdom"`
	Charisma     int64 `json:"charisma"`
}

var GetStatsReq struct {
	Id int64 `json:"id"`
}

type GetStatsRes struct {
	Strength     int64 `json:"strength"`
	Dexterity    int64 `json:"dexterity"`
	Constitution int64 `json:"constitution"`
	Intelligence int64 `json:"intelligence"`
	Wisdom       int64 `json:"wisdom"`
	Charisma     int64 `json:"charisma"`
}

type Repository interface {
	GetStatsById(ctx context.Context, id int64) (*Stats, error)
}

type Service interface {
	GetStatsById(c context.Context, id int64) (*GetStatsRes, error)
}
