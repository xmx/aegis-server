package repository

import "iter"

type Pages[E any, S ~[]*E] struct {
	Page    int64 `json:"page"    bson:"page"`
	Size    int64 `json:"size"    bson:"size"`
	Count   int64 `json:"count"   bson:"count"`
	Records S     `json:"records" bson:"records"`
}

func (pgs Pages[E, S]) All() iter.Seq2[int, *E] {
	return func(yield func(int, *E) bool) {
		for i, e := range pgs.Records {
			if !yield(i, e) {
				return
			}
		}
	}
}

type PageHelper interface {
	LPR(cnt int64) (page, limit, skip int64)
}

func NewPageHelper(page, size int64) PageHelper {
	return &pageHelper{page: page, size: size}
}

type pageHelper struct {
	page, size int64
}

// LPR Last Page Retention.
// 尾页保留，当请求页码超页后返回最后一页数据。
func (p *pageHelper) LPR(cnt int64) (page, limit, skip int64) {
	page, limit = p.page, p.size
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	} else if limit > 1000 {
		limit = 1000
	}

	if maximum := (cnt + limit - 1) / limit; maximum > 0 && maximum < page {
		page = maximum
	}
	skip = (page - 1) * limit

	return page, limit, skip
}
