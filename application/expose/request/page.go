package request

import (
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Sizes struct {
	Size int64 `json:"size" query:"size" form:"size" validate:"gte=0,lte=1000"`
}

func (s Sizes) Limit() int64 {
	if n := s.Size; n <= 0 {
		return 10
	} else if n > 1000 {
		return 1000
	} else {
		return n
	}
}

type Pages struct {
	Sizes
	Page int64 `json:"page" query:"page" form:"page" validate:"gte=0"`
}

type PageKeywords struct {
	Pages
	Keywords
}

type Keywords struct {
	Keyword string `json:"keyword" query:"keyword" form:"keyword" validate:"lte=255"`
}

func (k Keywords) Regexps(fields []string) bson.A {
	kw := strings.TrimSpace(k.Keyword)
	if len(fields) == 0 || kw == "" {
		return bson.A{}
	}
	kw = regexp.QuoteMeta(kw) // 防止用户注入特殊 regex 元字符
	reg := bson.Regex{Pattern: kw, Options: "i"}

	likes := make(bson.A, 0, len(fields))
	for _, f := range fields {
		likes = append(likes, bson.D{{Key: f, Value: reg}})
	}

	return likes
}
