package condition

import (
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type CondField struct {
	name    string
	comment string
	expr    field.Expr
}

func (f CondField) NameComment() (name, comment string) {
	return f.name, f.comment
}

type CondFields []*CondField

type WhereInput struct {
	Name    string
	Operate condOp
	Values  []string
}

type WhereInputs []*WhereInput

type OrderInput struct {
	Name string
	Desc bool
}

type OrderInputs []*OrderInput

type Cond struct {
	fileds  []*CondField
	nameMap map[string]*CondField
}

func (c Cond) OrderFields() CondFields {
	ret := make(CondFields, 0, 10)
	for _, f := range c.fileds {
		if _, ok := f.expr.(field.OrderExpr); ok {
			ret = append(ret, f)
		}
	}
	return ret
}

func (c Cond) Fields() CondFields {
	return c.fileds
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
		fd := c.nameMap[in.Name]
		if f, ok := fd.expr.(field.OrderExpr); ok {
			if in.Desc {
				exprs = append(exprs, f.Desc())
			} else {
				exprs = append(exprs, f.Asc())
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
	return nil
}

func (c Cond) parseWhereInput(input *WhereInput) gen.Condition {
	fd := c.nameMap[input.Name]
	if fd == nil {
		return nil
	}
	expr := fd.expr
	switch f := expr.(type) {
	case field.String:
		return c.fieldString(f, input)
	default:
		return nil
	}
}

func (c Cond) fieldString(f field.String, input *WhereInput) gen.Condition {
	values, op := input.Values, input.Operate
	size := len(values)
	if size == 0 {
		return nil
	}
	arg0 := values[0]

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
		return f.Regexp(arg0)
	case NotLike:
		return f.NotLike(arg0)
	case Regex:
		return f.Regexp(arg0)
	case NotRegex:
		return f.NotRegxp(arg0)
	case Between, NotBetween:
		if size < 2 {
			return nil
		}
		right := values[1]
		if op == Between {
			return f.Between(arg0, right)
		}
		return f.NotBetween(arg0, right)
	case In:
		return f.In(values...)
	case NotIn:
		return f.NotIn(values...)
	default:
		return nil
	}
}
