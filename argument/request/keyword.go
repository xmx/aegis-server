package request

import (
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type PageKeywords struct {
	Pages
	Keywords
}

type Keywords struct {
	Keyword string `json:"keyword" query:"keyword" form:"keyword" validate:"lte=100"`
}

func (k Keywords) Regexps(fields ...string) bson.A {
	regex := k.Regexp()
	if regex == nil || len(fields) == 0 {
		return nil
	}

	arr := make(bson.A, 0, len(fields))
	for _, field := range fields {
		arr = append(arr, bson.M{field: regex})
	}

	return arr
}

func (k Keywords) Regexp() *bson.Regex {
	if kw := strings.TrimSpace(k.Keyword); kw == "" {
		return nil
	} else {
		return &bson.Regex{Pattern: kw, Options: "i"}
	}
}
