package class

import "context"

type Class struct {
	Id        int64  `json:"id"`
	ClassName string `json:"className"`
}

type Repository interface {
	GetAllClasses(ctx context.Context) ([]Class, error)
}

type Service interface {
	GetAllClasses(ctx context.Context) ([]Class, error)
}
