package repository

import (
	"gorm.io/gen"
	"gorm.io/gorm"
)

type Page[T any] struct {
	Page    int64 `json:"page"`
	Size    int64 `json:"size"`
	Count   int64 `json:"count"`
	Records []T   `json:"records"`
}

type PageScope interface {
	Orm(count int64) func(db *gorm.DB) *gorm.DB
	Gen(count int64) func(dao gen.Dao) gen.Dao
	PageSize() (page, size int64)
}

func PageZero[T any](ps PageScope) *Page[T] {
	page, size := ps.PageSize()
	return &Page[T]{
		Page:    page,
		Size:    size,
		Records: []T{},
	}
}

func PageRecords[T any](ps PageScope, count int64, records []T) *Page[T] {
	page, size := ps.PageSize()
	return &Page[T]{
		Page:    page,
		Size:    size,
		Count:   count,
		Records: records,
	}
}
