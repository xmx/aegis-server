package request

import (
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type PageKeyword struct {
	Pages
	OptionalKeyword
}

type OptionalKeyword struct {
	Keyword string `json:"keyword" query:"keyword" form:"keyword" validate:"lte=100"`
}

func (o OptionalKeyword) LikeScope(fields ...field.String) func(gen.Dao) gen.Dao {
	size, kw := len(fields), o.Keyword
	if size == 0 || kw == "" {
		return func(dao gen.Dao) gen.Dao { return dao }
	}

	word := "%" + kw + "%"
	return func(dao gen.Dao) gen.Dao {
		conds := make([]gen.Condition, 0, size)
		for _, f := range fields {
			conds = append(conds, f.Like(word))
		}
		return dao.Or(conds...)
	}
}

func (o OptionalKeyword) Like() string {
	if kw := o.String(); kw != "" {
		return "%" + kw + "%"
	}
	return ""
}

func (o OptionalKeyword) LLike() string {
	if kw := o.String(); kw != "" {
		return "%" + kw
	}
	return ""
}

func (o OptionalKeyword) RLike() string {
	if kw := o.String(); kw != "" {
		return kw + "%"
	}
	return ""
}

func (o OptionalKeyword) String() string {
	return o.Keyword
}
