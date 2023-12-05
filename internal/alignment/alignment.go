package alignment

import "context"

type Alignment struct {
	Id            int64  `json:"id"`
	AlignmentName string `json:"alignmentName"`
}

type Repository interface {
	GetAllAlignments(ctx context.Context) ([]Alignment, error)
}

type Service interface {
	GetAllAlignments(ctx context.Context) ([]Alignment, error)
}
