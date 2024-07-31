package repository

import "github.com/xmx/aegis-server/datalayer/query"

type Repository[T any] interface {
	Query() *query.Query
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

//
//func (b *baseRepository[T]) Page(dao gen.Dao, conds []gen.Condition, scope PageScope) (*Page[T], error) {
//	stmt := dao.Where(conds...)
//	count, err := stmt.Count()
//	if err != nil {
//		return nil, err
//	}
//	if count == 0 {
//		return b.empty(scope), nil
//	}
//	dats, err := dao.Scopes(scope.Gen(count)).Find()
//	if err != nil {
//		return nil, err
//	}
//	b.record(scope, count, dats)
//	b.qry.ConfigCertificate.WithContext()
//}
//
//func (*baseRepository[T]) empty(p PageScope) *Page[T] {
//	_, size := p.PageSize()
//	return &Page[T]{
//		Page:    1,
//		Size:    size,
//		Records: []T{},
//	}
//}
//
//func (*baseRepository[T]) record(p PageScope, count int64, dats []T) *Page[T] {
//	page, size := p.PageSize()
//	return &Page[T]{
//		Page:    page,
//		Size:    size,
//		Count:   count,
//		Records: dats,
//	}
//}
