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
	name      string
	comment   string
	operators []Operator
	expr      field.Expr
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

type WhereInputs struct {
	Or     bool
	Inputs []*WhereInput
}

type OrderInput struct {
	Name string
	Desc bool
}

type OrderInputs struct {
	Inputs []*OrderInput
}

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

func (c Cond) Scope(whereInputs *WhereInputs, orderInputs *OrderInputs) func(gen.Dao) gen.Dao {
	return func(dao gen.Dao) gen.Dao {
		wheres := c.parseWhereInputs(whereInputs)
		orders := c.parseOrderInputs(orderInputs)
		return dao.Where(wheres...).Order(orders...)
	}
}

func (c Cond) Order(inputs *OrderInputs) func(gen.Dao) gen.Dao {
	return func(dao gen.Dao) gen.Dao {
		exprs := c.parseOrderInputs(inputs)
		return dao.Order(exprs...)
	}
}

func (c Cond) Where(inputs *WhereInputs) func(gen.Dao) gen.Dao {
	return func(dao gen.Dao) gen.Dao {
		exprs := c.parseWhereInputs(inputs)
		return dao.Where(exprs...)
	}
}

func (c Cond) parseOrderInputs(in *OrderInputs) []field.Expr {
	if in == nil {
		return nil
	}
	inputs := in.Inputs
	exprs := make([]field.Expr, 0, len(inputs))
	for _, v := range inputs {
		if in == nil {
			continue
		}
		if fd := c.ordersNameMap[v.Name]; fd != nil {
			exp := fd.expr
			if v.Desc {
				exprs = append(exprs, exp.Desc())
			} else {
				exprs = append(exprs, exp.Asc())
			}
		}
	}
	return exprs
}

func (c Cond) parseWhereInputs(in *WhereInputs) []gen.Condition {
	if in == nil {
		return nil
	}
	or, inputs := in.Or, in.Inputs
	size := len(inputs)
	exprs := make([]field.Expr, 0, size)
	conds := make([]gen.Condition, 0, size)
	for _, val := range inputs {
		expr := c.parseWhereInput(val)
		exprs = append(exprs, expr)
		conds = append(conds, expr)
	}
	if or {
		return []gen.Condition{field.Or(exprs...)}
	}
	return conds
}

func (c Cond) parseWhereInput(input *WhereInput) field.Expr {
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
	case field.Int64:
		return c.fieldInt64(f, input)
	case field.Field:
		return c.fieldField(f, input)
	default:
		return nil
	}
}

func (c Cond) fieldString(f field.String, input *WhereInput) field.Expr {
	values, op := stringValues(input.Values), input.Operate
	arg0, ok := values.stringN(0)
	if !ok {
		return nil
	}
	if !ok || arg0 == "" {
		switch op {
		case Eq, Neq:
		default:
			return nil
		}
	}

	switch op {
	case Eq:
		if !ok {
			return f.IsNull()
		}
		return f.Eq(arg0)
	case Neq:
		if !ok {
			return f.IsNotNull()
		}
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
		arg1, exist := values.stringN(1)
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

func (c Cond) fieldInt(f field.Int, input *WhereInput) field.Expr {
	values, op := stringValues(input.Values), input.Operate
	arg0, ok := values.intN(0)
	if !ok {
		switch op {
		case Eq, Neq:
		default:
			return nil
		}
	}

	switch op {
	case Eq:
		if ok {
			return f.Eq(arg0)
		}
		return f.IsNull()
	case Neq:
		if ok {
			return f.IsNotNull()
		}
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

func (c Cond) fieldInt64(f field.Int64, input *WhereInput) field.Expr {
	values, op := stringValues(input.Values), input.Operate
	arg0, ok := values.int64N(0)
	if !ok && op != Eq && op != Neq {
		return nil
	}

	switch op {
	case Eq:
		if ok {
			return f.Eq(arg0)
		}
		return f.IsNull()
	case Neq:
		if ok {
			return f.Neq(arg0)
		}
		return f.IsNotNull()
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
		arg1, exist := values.int64N(1)
		if !exist {
			return nil
		}
		if op == Between {
			return f.Between(arg0, arg1)
		}
		return f.NotBetween(arg0, arg1)
	case In, NotIn:
		args := values.int64s()
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

func (c Cond) fieldTime(f field.Time, input *WhereInput) field.Expr {
	values, op := stringValues(input.Values), input.Operate
	arg0, ok := values.timeN(0)
	if !ok && op != Eq && op != Neq {
		return nil
	}

	switch op {
	case Eq:
		if ok {
			return f.Eq(arg0)
		}
		return f.IsNull()
	case Neq:
		if ok {
			return f.Neq(arg0)
		}
		return f.IsNotNull()
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

func (c Cond) fieldBool(f field.Bool, input *WhereInput) field.Expr {
	values, op := stringValues(input.Values), input.Operate
	arg0, ok := values.boolN(0)
	if !ok && op != Eq && op != Neq {
		return nil
	}
	switch op {
	case Eq:
		if ok {
			return f.Is(arg0)
		}
		return f.IsNull()
	case Neq:
		if ok {
			return f.Is(!arg0)
		}
		return f.IsNotNull()
	default:
		return nil
	}
}

func (c Cond) fieldField(f field.Field, input *WhereInput) field.Expr {
	values, op := stringValues(input.Values), input.Operate
	arg0, ok := values.valueN(0)
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
	case In, NotIn:
		vals := values.values()
		if len(vals) == 0 {
			return nil
		}
		return f.In(vals...)
	default:
		return nil
	}
}
