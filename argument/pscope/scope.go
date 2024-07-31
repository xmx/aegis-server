package pscope

import (
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/argument/response"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func Zero[T any](p request.Page) *response.Page[T] {
	size := p.Size
	if size <= 0 {
		size = 10
	}
	return &response.Page[T]{
		Page:    1,
		Size:    p.Size,
		Records: []T{},
	}
}

func With[T any](p request.Page, cnt int64) Scope[T] {
	return Page[T](p.Page, p.Size, cnt)
}

func Page[T any](page int64, size int64, cnt int64) Scope[T] {
	if page < 1 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	return &pageScope[T]{
		page:  page,
		size:  size,
		count: cnt,
	}
}

type Scope[T any] interface {
	Orm(db *gorm.DB) *gorm.DB
	Gen(dao gen.Dao) gen.Dao
	Records(ts []T) *response.Page[T]
}

type pageScope[T any] struct {
	page  int64
	size  int64
	count int64
}

func (p *pageScope[T]) Orm(db *gorm.DB) *gorm.DB {
	offset, limit := p.LPR()
	return db.Offset(offset).Limit(limit)
}

func (p *pageScope[T]) Gen(dao gen.Dao) gen.Dao {
	offset, limit := p.LPR()
	return dao.Offset(offset).Limit(limit)
}

func (p *pageScope[T]) LPR() (offset int, limit int) {
	if p.count >= 0 {
		pg := (p.count + p.size - 1) / p.size
		if pg == 0 {
			p.page = 1
		} else if pg < p.page {
			p.page = pg
		}
	}
	offset, limit = int((p.page-1)*p.size), int(p.size)

	return
}

func (p *pageScope[T]) Records(ts []T) *response.Page[T] {
	ret := &response.Page[T]{
		Page:    p.page,
		Size:    p.size,
		Count:   p.count,
		Records: ts,
	}
	if ts == nil {
		ret.Records = []T{}
	}

	return ret
}
