package response

import "github.com/xmx/aegis-server/datalayer/dynsql"

type CondEnums []*CondEnum

type CondEnum struct {
	Value any    `json:"value"` // 枚举值
	Name  string `json:"name"`  // 枚举说明
}

type CondOrders []*CondOrder

type CondOrder struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CondOperators []*CondOperator

type CondOperator struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CondWheres []*CondWhere

type CondWhere struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Type      string        `json:"type"`
	Enums     CondEnums     `json:"enums,omitempty,omitzero"`
	Operators CondOperators `json:"operators"`
}

type Cond struct {
	Wheres CondWheres `json:"wheres,omitempty,omitzero"`
	Orders CondOrders `json:"orders,omitempty,omitzero"`
}

func ReadCond(tbl *dynsql.Table) *Cond {
	cd := new(Cond)
	wheres, orders := tbl.RawColumns()
	for _, w := range wheres {
		enums := make(CondEnums, 0, len(w.Enums))
		for _, e := range w.Enums {
			ce := &CondEnum{Value: e.Value, Name: e.Name}
			enums = append(enums, ce)
		}
		operators := make(CondOperators, 0, len(w.Operators))
		for _, o := range w.Operators {
			id, name := o.Info()
			operators = append(operators, &CondOperator{ID: id, Name: name})
		}

		cw := &CondWhere{
			ID:        w.ID,
			Name:      w.Name,
			Type:      string(w.Type),
			Enums:     enums,
			Operators: operators,
		}
		cd.Wheres = append(cd.Wheres, cw)
	}
	for _, o := range orders {
		co := &CondOrder{ID: o.ID, Name: o.Name}
		cd.Orders = append(cd.Orders, co)
	}

	return cd
}
