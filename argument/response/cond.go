package response

import (
	"github.com/xmx/aegis-server/datalayer/condition"
	"github.com/xmx/aegis-server/datalayer/dynsql"
)

type Cond struct {
	Wheres []*CondWhere `json:"wheres,omitempty"`
	Orders []*CondOrder `json:"orders,omitempty"`
}

type CondOrder struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type CondWhere struct {
	Name      string   `json:"name"`
	Comment   string   `json:"comment"`
	Operators []string `json:"operators"`
}

func ReadCond(c *condition.Cond) *Cond {
	ret := new(Cond)
	if c == nil {
		return ret
	}

	for _, f := range c.OrderFields() {
		name, comment := f.NameComment()
		data := &CondOrder{Name: name, Comment: comment}
		ret.Orders = append(ret.Orders, data)
	}
	for _, f := range c.WhereFields() {
		name, comment := f.NameComment()
		data := &CondWhere{Name: name, Comment: comment}
		ops := f.Operators()
		for _, op := range ops {
			data.Operators = append(data.Operators, op.String())
		}

		ret.Wheres = append(ret.Wheres, data)
	}

	return ret
}

type CondEnums []*CondEnum

type CondEnum struct {
	Value any    `json:"value"` // 枚举值
	Name  string `json:"name"`  // 枚举说明
}

type CondOrders []*CondOrder1

type CondOrder1 struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CondOperators []*CondOperator

type CondOperator struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CondWheres []*CondWhere1

type CondWhere1 struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Type      string        `json:"type"`
	Enums     CondEnums     `json:"enums,omitempty,omitzero"`
	Operators CondOperators `json:"operators"`
}

type Cond1 struct {
	Wheres CondWheres `json:"wheres,omitempty,omitzero"`
	Orders CondOrders `json:"orders,omitempty,omitzero"`
}

func ReadCond1(tbl *dynsql.Table) *Cond1 {
	cd := new(Cond1)
	wheres, orders := tbl.RawColumns()
	for _, w := range wheres {
		enums := make(CondEnums, 0, len(w.Enums))
		for _, e := range w.Enums {
			ce := &CondEnum{Value: e.Value, Name: e.Name}
			enums = append(enums, ce)
		}
		operators := make(CondOperators, 0, len(w.Operators))
		for _, o := range w.Operators {
			id, name := o.OpInfo()
			operators = append(operators, &CondOperator{ID: id, Name: name})
		}

		cw := &CondWhere1{
			ID:        w.ID,
			Name:      w.Name,
			Type:      string(w.Type),
			Enums:     enums,
			Operators: operators,
		}
		cd.Wheres = append(cd.Wheres, cw)
	}
	for _, o := range orders {
		co := &CondOrder1{ID: o.ID, Name: o.Name}
		cd.Orders = append(cd.Orders, co)
	}

	return cd
}
