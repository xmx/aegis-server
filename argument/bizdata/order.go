package bizdata

import (
	"encoding/json"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

type FieldDesc struct {
	Field string `json:"field"` // 字段名
	Desc  bool   `json:"desc"`  // 是否倒序
}

func (fd *FieldDesc) UnmarshalBind(param string) error {
	return json.Unmarshal([]byte(param), fd)
}

type FieldComment struct {
	Field   string `json:"field"`
	Comment string `json:"comment"`
}

type SearchOrders struct {
	fields  []*FieldComment
	indices map[string]field.OrderExpr
}

func (so *SearchOrders) Add(f field.OrderExpr, comment string) *SearchOrders {
	so.initial()

	name := string(f.ColumnName())
	if name != "" {
		if comment == "" {
			comment = name
		}
		fc := &FieldComment{Field: name, Comment: comment}
		so.fields = append(so.fields, fc)
		so.indices[name] = f
	}

	return so
}

func (so *SearchOrders) Fields() []*FieldComment {
	return so.fields
}

func (so *SearchOrders) Exprs(fos []*FieldDesc) ([]field.Expr, []*FieldDesc) {
	so.initial()

	size := len(fos)
	mismatch := make([]*FieldDesc, 0, size)
	orders := make([]field.Expr, 0, size)
	for _, fo := range fos {
		expr, ok := so.indices[fo.Field]
		if !ok {
			mismatch = append(mismatch, fo)
			continue
		}

		ord := expr.Asc()
		if fo.Desc {
			ord = expr.Desc()
		}
		orders = append(orders, ord)
	}
	if len(orders) == 0 {
		return nil, mismatch
	}

	return orders, mismatch
}

func (so *SearchOrders) Scope(fos []*FieldDesc) func(dao gen.Dao) gen.Dao {
	size := len(fos)
	if so == nil || size == 0 {
		return func(dao gen.Dao) gen.Dao {
			return dao
		}
	}

	so.initial()
	orders := make([]field.Expr, 0, size)
	for _, fo := range fos {
		expr, ok := so.indices[fo.Field]
		if !ok {
			continue
		}
		ord := expr.Asc()
		if fo.Desc {
			ord = expr.Desc()
		}
		orders = append(orders, ord)
	}
	return func(dao gen.Dao) gen.Dao {
		if len(orders) == 0 {
			return dao
		}

		return dao.Order(orders...)
	}
}

func (so *SearchOrders) initial() {
	if so.indices == nil {
		so.indices = make(map[string]field.OrderExpr, 8)
	}
}
