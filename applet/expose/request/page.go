package request

import (
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Pages struct {
	Page int64 `json:"page" query:"page" form:"page" validate:"gte=0"`
	Size int64 `json:"size" query:"size" form:"size" validate:"gte=0,lte=1000"`
}

type PageKeywords struct {
	Pages
	Keywords
}

type Keywords struct {
	Keyword string `json:"keyword" query:"keyword" form:"keyword" validate:"lte=255"`
}

func (k Keywords) Regexps(fields []string) bson.D {
	kw := strings.TrimSpace(k.Keyword)
	if len(fields) == 0 || kw == "" {
		return bson.D{}
	}
	kw = regexp.QuoteMeta(kw) // 防止用户注入特殊 regex 元字符
	reg := bson.Regex{Pattern: kw, Options: "i"}

	likes := make(bson.A, 0, len(fields))
	for _, f := range fields {
		likes = append(likes, bson.D{{Key: f, Value: reg}})
	}

	return bson.D{{Key: "$or", Value: likes}}
}
