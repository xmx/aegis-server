package repository

func NewPages[T any](page, size int64) *Pages[T] {
	return &Pages[T]{
		Page:    page,
		Size:    size,
		Records: []*T{},
	}
}

type Pages[T any] struct {
	Page    int64 `json:"page"`
	Size    int64 `json:"size"`
	Total   int64 `json:"total"`
	Records []*T  `json:"records"`
}

func (ps *Pages[T]) Skip(cnt int64) int64 {
	ps.Total = cnt
	if ps.Page <= 0 {
		ps.Page = 1
	}
	if ps.Size <= 0 {
		ps.Size = 10
	}
	page, limit := ps.Page, ps.Size
	if last := (cnt + limit - 1) / limit; last > 0 && last < page {
		page = last
	}

	return (page - 1) * limit
}
