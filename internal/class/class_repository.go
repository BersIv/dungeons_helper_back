package class

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

func (r *repository) GetAllClasses(ctx context.Context) ([]Class, error) {
	var classes []Class

	query := "SELECT id, className FROM class"
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
		var class Class
		err := rows.Scan(&class.Id, &class.ClassName)
		if err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return classes, nil
}
