package stats

import (
	"context"
	"dungeons_helper/db"
)

type repository struct {
	db db.DatabaseTX
}

func NewRepository(db db.DatabaseTX) Repository {
	return &repository{db: db}
}

func (r *repository) GetStatsById(ctx context.Context, id int64) (*Stats, error) {
	stats := Stats{}
	query := "SELECT strength, dexterity, constitution, intelligence, wisdom, charisma FROM stats WHERE id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&stats.Strength, &stats.Dexterity,
		&stats.Constitution, &stats.Intelligence, &stats.Wisdom, &stats.Charisma)
	if err != nil {
		return &Stats{}, err
	}

	return &stats, nil
}
