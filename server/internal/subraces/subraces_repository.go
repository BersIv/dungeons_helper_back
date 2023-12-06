package subraces

import (
	"context"
	"database/sql"
	"dungeons_helper/db"
)

type repository struct {
	db db.DatabaseTX
}

func NewRepository(db db.DatabaseTX) Repository {
	return &repository{db: db}
}

func (r *repository) GetAllSubraces(ctx context.Context) ([]Subraces, error) {
	var subraces []Subraces

	query := "SELECT id, subraceName, idStats FROM subrace"
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
		var subrace Subraces
		err := rows.Scan(&subrace.Id, &subrace.SubraceName, &subrace.IdStats)
		if err != nil {
			return nil, err
		}
		subraces = append(subraces, subrace)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return subraces, nil
}
