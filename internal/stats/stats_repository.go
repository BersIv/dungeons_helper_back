package stats

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
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

func (r *repository) CreateStats(ctx context.Context, stats *Stats) (*Stats, error) {
	query := "INSERT INTO stats(strength, dexterity, constitution, intelligence, wisdom, charisma) VALUES (?, ?, ?, ?, ?)"
	result, err := r.db.ExecContext(ctx, query, stats.Strength, stats.Dexterity, stats.Constitution,
		stats.Intelligence, stats.Wisdom, stats.Charisma)
	if err != nil {
		return &Stats{}, err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return &Stats{}, err
	}
	stats.Id = int64(lastInsertID)

	return stats, nil
}
