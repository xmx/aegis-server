package gormcond

import (
	"slices"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

type ColumnComment struct {
	Column  string
	Comment string
}

type ColumnDesc struct {
	Column string
	Desc   bool
}

func NewOrder() *Order {
	return &Order{
		mapping: make(map[string]field.OrderExpr, 8),
		columns: make([]*ColumnComment, 0, 8),
	}
}

type Order struct {
	mapping map[string]field.OrderExpr
	columns []*ColumnComment
}

func (o *Order) Add(f field.OrderExpr, comment string) *Order {
	o.initial()
	if f == nil {
		return o
	}
	name := f.ColumnName().String()
	if name == "" {
		return o
	}
	if comment == "" {
		comment = name
	}
	_, exist := o.mapping[name]
	o.mapping[name] = f
	o.columns = append(o.columns, &ColumnComment{Column: name, Comment: comment})
	if exist {
		var found bool
		o.columns = slices.DeleteFunc(o.columns, func(cc *ColumnComment) bool {
			if !found && cc.Column == name {
				found = true
				return true
			}
			return false
		})
	}

	return o
}

func (o *Order) Columns() []*ColumnComment {
	return o.columns
}

func (o *Order) Scope(cds []*ColumnDesc) func(gen.Dao) gen.Dao {
	return func(dao gen.Dao) gen.Dao {
		size := len(cds)
		if size == 0 || len(o.mapping) == 0 {
			return dao
		}

		columns := make([]field.Expr, 0, size)
		for _, fo := range cds {
			col := fo.Column
			fc := o.mapping[col]
			if fc == nil {
				continue
			}
			if fo.Desc {
				columns = append(columns, fc.Desc())
			} else {
				columns = append(columns, fc.Asc())
			}
		}
		if len(columns) == 0 {
			return dao
		}

		return dao.Order(columns...)
	}
}

func (o *Order) initial() {
	if o.mapping == nil {
		o.mapping = make(map[string]field.OrderExpr, 8)
	}
	if o.columns == nil {
		o.columns = make([]*ColumnComment, 0, 8)
	}
}
