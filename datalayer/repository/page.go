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
	Orm(cnt int64) func(db *gorm.DB) *gorm.DB
	Gen(cnt int64) func(dao gen.Dao) gen.Dao
	PageSize() (page, size int64)
}
