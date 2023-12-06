package skills

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

func (r *repository) GetAllSkills(ctx context.Context) ([]Skills, error) {
	var skills []Skills

	query := "SELECT id, skillName FROM skills"
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
		var skill Skills
		err := rows.Scan(&skill.Id, &skill.SkillName)
		if err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return skills, nil
}
