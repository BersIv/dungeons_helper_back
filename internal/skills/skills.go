package skills

import "context"

type Skills struct {
	Id        int64  `json:"id"`
	SkillName string `json:"skillName"`
}

type Repository interface {
	GetAllSkills(ctx context.Context) ([]Skills, error)
}

type Service interface {
	GetAllSkills(ctx context.Context) ([]Skills, error)
}
