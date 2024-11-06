package pagination

import "gorm.io/gen"

type Pager[E any] interface {
	Scope(cnt int64) func(gen.Dao) gen.Dao
	Empty() *Result[E]
	Result(records []E) *Result[E]
}

func NewPager[E any](page, size int64) Pager[E] {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	} else if size >= 1000 {
		size = 1000
	}

	return &pagination[E]{
		page:  page,
		size:  size,
		count: 0,
	}
}

type pagination[E any] struct {
	page  int64
	size  int64
	count int64
}

func (p *pagination[E]) Empty() *Result[E] {
	return &Result[E]{
		Page:    p.page,
		Size:    p.size,
		Records: []E{},
	}
}

func (p *pagination[E]) Result(records []E) *Result[E] {
	return &Result[E]{
		Page:    p.page,
		Size:    p.size,
		Count:   p.count,
		Records: records,
	}
}

func (p *pagination[E]) Scope(cnt int64) func(gen.Dao) gen.Dao {
	offset, limit := p.calculate(cnt)
	return func(dao gen.Dao) gen.Dao {
		return dao.Offset(offset).Limit(limit)
	}
}

func (p *pagination[E]) calculate(cnt int64) (offset, limit int) {
	if cnt >= 0 {
		p.count = cnt
		page := (cnt + p.size - 1) / p.size
		if page <= 0 {
			p.page = 1
		} else if page < p.page {
			p.page = page
		}
	}
	offset = int((p.page - 1) * p.size)
	limit = int(p.size)

	return
}
