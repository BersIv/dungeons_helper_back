package races

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

func (r *repository) GetAllRaces(ctx context.Context) ([]Races, error) {
	var races []Races

	query := "SELECT id, raceName FROM races"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	for rows.Next() {
		var race Races
		err := rows.Scan(&race.Id, &race.RaceName)
		if err != nil {
			return nil, err
		}
		races = append(races, race)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return races, nil
}
