package pscope

import (
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/datalayer/repository"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func With(p *request.Page) repository.PageScope {
	if p == nil {
		return Page(1, 10)
	}

	return From(*p)
}

func From(p request.Page) repository.PageScope {
	return Page(p.Page, p.Size)
}

func Page(page, size int64) repository.PageScope {
	if page < 1 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	return &pageScope{
		page: page,
		size: size,
	}
}

type pageScope struct {
	page  int64
	size  int64
	count int64
}

func (ps *pageScope) Orm(count int64) func(db *gorm.DB) *gorm.DB {
	offset, limit := ps.LPR(count)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit)
	}
}

func (ps *pageScope) Gen(count int64) func(dao gen.Dao) gen.Dao {
	offset, limit := ps.LPR(count)

	return func(dao gen.Dao) gen.Dao {
		return dao.Offset(offset).Limit(limit)
	}
}

func (ps *pageScope) PageSize() (page, size int64) {
	return ps.size, ps.size
}

func (ps *pageScope) LPR(count int64) (offset int, limit int) {
	if count >= 0 {
		ps.count = count
		pg := (count + ps.size - 1) / ps.size
		if pg == 0 {
			ps.page = 1
		} else if pg < ps.page {
			ps.page = pg
		}
	}
	offset, limit = int((ps.page-1)*ps.size), int(ps.size)

	return
}
