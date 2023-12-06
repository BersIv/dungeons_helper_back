package alignment

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

func (r *repository) GetAllAlignments(ctx context.Context) ([]Alignment, error) {
	var alignments []Alignment

	query := "SELECT id, alignmentName FROM alignment"
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
		var alignment Alignment
		err := rows.Scan(&alignment.Id, &alignment.AlignmentName)
		if err != nil {
			return nil, err
		}
		alignments = append(alignments, alignment)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return alignments, nil
}
