package condition

import (
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type OrderField struct {
	name    string
	comment string
	expr    field.OrderExpr
}

func (f OrderField) NameComment() (name, comment string) {
	return f.name, f.comment
}

type OrderFields []*OrderField

type WhereField struct {
	name    string
	comment string
	expr    field.Expr
}

func (f WhereField) NameComment() (name, comment string) {
	return f.name, f.comment
}

type WhereFields []*WhereField

type WhereInput struct {
	Name    string
	Operate Operator
	Values  []string
}

type WhereInputs []*WhereInput

type OrderInput struct {
	Name string
	Desc bool
}

type OrderInputs []*OrderInput

type Cond struct {
	orders        []*OrderField
	ordersNameMap map[string]*OrderField
	wheres        []*WhereField
	wheresNameMap map[string]*WhereField
}

func (c Cond) OrderFields() OrderFields {
	return c.orders
}

func (c Cond) WhereFields() WhereFields {
	return c.wheres
}

func (c Cond) Order(inputs OrderInputs) func(gen.Dao) gen.Dao {
	return func(dao gen.Dao) gen.Dao {
		exprs := c.parseOrderInputs(inputs)
		return dao.Order(exprs...)
	}
}

func (c Cond) Where(inputs WhereInputs) func(gen.Dao) gen.Dao {
	return func(dao gen.Dao) gen.Dao {
		exprs := c.parseWhereInputs(inputs)
		return dao.Where(exprs...)
	}
}

func (c Cond) parseOrderInputs(inputs OrderInputs) []field.Expr {
	exprs := make([]field.Expr, 0, len(inputs))
	for _, in := range inputs {
		if in == nil {
			continue
		}
		if fd := c.ordersNameMap[in.Name]; fd != nil {
			exp := fd.expr
			if in.Desc {
				exprs = append(exprs, exp.Desc())
			} else {
				exprs = append(exprs, exp.Asc())
			}
		}
	}
	return exprs
}

func (c Cond) parseWhereInputs(inputs WhereInputs) []gen.Condition {
	exprs := make([]gen.Condition, 0, len(inputs))
	for _, in := range inputs {
		if input := c.parseWhereInput(in); input != nil {
			exprs = append(exprs, input)
		}
	}
	return exprs
}

func (c Cond) parseWhereInput(input *WhereInput) gen.Condition {
	fd := c.wheresNameMap[input.Name]
	if fd == nil {
		return nil
	}
	expr := fd.expr
	switch f := expr.(type) {
	case field.String:
		return c.fieldString(f, input)
	case field.Bool:
		return c.fieldBool(f, input)
	case field.Time:
		return c.fieldTime(f, input)
	case field.Int:
		return c.fieldInt(f, input)
	default:
		return nil
	}
}

func (c Cond) fieldString(f field.String, input *WhereInput) gen.Condition {
	values, op := stringValues(input.Values), input.Operate
	arg0, ok := values.getN(0)
	if !ok {
		return nil
	}

	switch op {
	case Eq:
		return f.Eq(arg0)
	case Neq:
		return f.Neq(arg0)
	case Gt:
		return f.Gt(arg0)
	case Gte:
		return f.Gte(arg0)
	case Lt:
		return f.Lt(arg0)
	case Lte:
		return f.Lte(arg0)
	case Like:
		return f.Like(arg0)
	case NotLike:
		return f.NotLike(arg0)
	case Regex:
		return f.Regexp(arg0)
	case NotRegex:
		return f.NotRegxp(arg0)
	case Between, NotBetween:
		arg1, exist := values.getN(1)
		if !exist {
			return nil
		}
		if op == Between {
			return f.Between(arg0, arg1)
		}
		return f.NotBetween(arg0, arg1)
	case In:
		return f.In(values...)
	case NotIn:
		return f.NotIn(values...)
	default:
		return nil
	}
}

func (c Cond) fieldInt(f field.Int, input *WhereInput) gen.Condition {
	values, op := stringValues(input.Values), input.Operate
	arg0, ok := values.intN(0)
	if !ok {
		return nil
	}

	switch op {
	case Eq:
		return f.Eq(arg0)
	case Neq:
		return f.Neq(arg0)
	case Gt:
		return f.Gt(arg0)
	case Gte:
		return f.Gte(arg0)
	case Lt:
		return f.Lt(arg0)
	case Lte:
		return f.Lte(arg0)
	case Like:
		return f.Like(arg0)
	case NotLike:
		return f.NotLike(arg0)
	case Between, NotBetween:
		arg1, exist := values.intN(1)
		if !exist {
			return nil
		}
		if op == Between {
			return f.Between(arg0, arg1)
		}
		return f.NotBetween(arg0, arg1)
	case In, NotIn:
		args := values.ints()
		if len(args) == 0 {
			return nil
		}
		if op == In {
			return f.In(args...)
		}
		return f.NotIn(args...)
	default:
		return nil
	}
}

func (c Cond) fieldTime(f field.Time, input *WhereInput) gen.Condition {
	values, op := stringValues(input.Values), input.Operate
	arg0, ok := values.timeN(0)
	if !ok {
		return nil
	}

	switch op {
	case Eq:
		return f.Eq(arg0)
	case Neq:
		return f.Neq(arg0)
	case Gt:
		return f.Gt(arg0)
	case Gte:
		return f.Gte(arg0)
	case Lt:
		return f.Lt(arg0)
	case Lte:
		return f.Lte(arg0)
	case Between, NotBetween:
		arg1, exist := values.timeN(1)
		if !exist {
			return nil
		}
		if op == Between {
			return f.Between(arg0, arg1)
		}
		return f.NotBetween(arg0, arg1)
	case In, NotIn:
		args := values.times()
		if len(args) == 0 {
			return nil
		}
		if op == In {
			return f.In(args...)
		}
		return f.NotIn(args...)
	default:
		return nil
	}
}

func (c Cond) fieldBool(f field.Bool, input *WhereInput) gen.Condition {
	values, op := stringValues(input.Values), input.Operate
	arg0, ok := values.boolN(0)
	if !ok {
		return nil
	}
	switch op {
	case Eq:
		return f.Is(arg0)
	default:
		return nil
	}
}
