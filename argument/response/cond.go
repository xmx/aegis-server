package response

import "github.com/xmx/aegis-server/datalayer/condition"

type Cond struct {
	Wheres []*CondOrder `json:"wheres,omitempty"`
	Orders []*CondOrder `json:"orders,omitempty"`
}

type CondOrder struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
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
		data := &CondOrder{Name: name, Comment: comment}
		ret.Wheres = append(ret.Wheres, data)
	}

	return ret
}
