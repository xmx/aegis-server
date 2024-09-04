package repository

import (
	"database/sql"
	"time"

	"github.com/xmx/aegis-server/datalayer/query"
)

type Repository[T any] interface {
	Query() *query.Query

	emptyRecords(ps PageScope) *Page[T]

	withRecords(ps PageScope, count int64, records []T) *Page[T]
}

func Base[T any](qry *query.Query) Repository[T] {
	return &baseRepository[T]{qry: qry}
}

type baseRepository[T any] struct {
	qry *query.Query
}

func (b *baseRepository[T]) Query() *query.Query {
	return b.qry
}

func (*baseRepository[T]) emptyRecords(ps PageScope) *Page[T] {
	page, size := ps.PageSize()
	return &Page[T]{
		Page:    page,
		Size:    size,
		Records: []T{},
	}
}

func (*baseRepository[T]) withRecords(ps PageScope, count int64, records []T) *Page[T] {
	page, size := ps.PageSize()
	return &Page[T]{
		Page:    page,
		Size:    size,
		Count:   count,
		Records: records,
	}
}

type name struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
	ExpiredAt sql.NullTime
}
