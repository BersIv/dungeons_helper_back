package races

import (
	"context"
	"database/sql"
	"dungeons_helper_server/db"
)

type repository struct {
	db db.DatabaseTX
}

func NewRepository(db db.DatabaseTX) Repository {
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
