package dynsql

import (
	"errors"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

type Table struct {
	wheres    Wheres
	orders    Orders
	whereMaps map[string]*Where
	orderMaps map[string]*Order
}

func (t Table) RawColumns() (Wheres, Orders) {
	return t.wheres, t.orders
}

func (t Table) ScopeWheres(inputs WhereInputs) func(gen.Dao) gen.Dao {
	return func(dao gen.Dao) gen.Dao {
		if build, err := t.BuildWheres(inputs); err != nil {
			_ = dao.AddError(err)
			return dao
		} else {
			return dao.Where(build.Conds...)
		}
	}
}

func (t Table) ScopeOrders(inputs OrderInputs) func(gen.Dao) gen.Dao {
	return func(dao gen.Dao) gen.Dao {
		if build, err := t.BuildOrders(inputs); err != nil {
			_ = dao.AddError(err)
			return dao
		} else {
			return dao.Order(build.Exprs...)
		}
	}
}

func (t Table) BuildWheres(inputs WhereInputs) (*WhereBuild, error) {
	build := new(WhereBuild)
	for _, input := range inputs {
		cond, expr, err := t.buildWhere(input)
		if err != nil {
			return nil, err
		}
		build.Conds = append(build.Conds, cond)
		build.Exprs = append(build.Exprs, cond)
		build.Hits = append(build.Hits, expr)
	}

	return build, nil
}

func (t Table) BuildOrders(inputs OrderInputs) (*OrderBuild, error) {
	build := new(OrderBuild)
	for _, input := range inputs {
		cond, expr, err := t.buildOrder(input)
		if err != nil {
			return nil, err
		}
		build.Exprs = append(build.Exprs, cond)
		build.Hits = append(build.Hits, expr)
	}

	return build, nil
}

func (t Table) buildOrder(input *OrderInput) (field.Expr, field.Expr, error) {
	order := t.orderMaps[input.ID]
	if order == nil {
		return nil, nil, errors.New("order not found for " + input.ID)
	}

	expr := order.Expr
	if input.Desc {
		return expr.Desc(), expr, nil
	}

	return expr.Asc(), expr, nil
}

func (t Table) buildWhere(input *WhereInput) (field.Expr, field.Expr, error) {
	where := t.whereMaps[input.ID]
	if where == nil {
		return nil, nil, errors.New(input.ID + " is not exist")
	}
	id, _ := input.Operator.OpInfo()
	op := where.opMaps[id]
	if op == nil {
		return nil, nil, errors.New("operator " + id + " is not exist")
	}

	expr, values := where.Expr, stringValues(input.Values)
	vsize := len(values)
	if vsize == 0 && op != Null {
		return nil, nil, errors.New("参数必须存在")
	}
	if (op == Between || op == NotBetween) && vsize < 2 {
		return nil, nil, errors.New("该操作参数必须存在2个参数")
	}

	var err error
	var cond field.Expr
	switch f := expr.(type) {
	case field.String:
		cond, err = t.buildString(f, op, values)
	case field.Int:
	case field.Float64:
	case field.Bool:
	case field.Time:

	}
	if err != nil {
		return nil, nil, err
	}

	return cond, expr, nil
}

func (t Table) buildString(f field.String, op Operator, vals stringValues) (field.Expr, error) {
	arg0, _ := vals.string0()
	switch op {
	case Eq:
		return f.Eq(arg0), nil
	case Neq:
		return f.Neq(arg0), nil
	case Gt:
		return f.Gt(arg0), nil
	case Gte:
		return f.Gte(arg0), nil
	case Lt:
		return f.Lt(arg0), nil
	case Lte:
		return f.Lte(arg0), nil
	case Like:
	case NotLike:
	case Between:
		arg1, _ := vals.stringN(1)
		return f.Between(arg0, arg1), nil
	case NotBetween:
		arg1, _ := vals.stringN(1)
		return f.NotBetween(arg0, arg1), nil
	case In:
		return f.In(vals...), nil
	case NotIn:
		return f.NotIn(vals...), nil
	case Null:
		nullable, exists := vals.bool0()
		if exists && !nullable {
			return f.IsNotNull(), nil
		}
		return f.IsNull(), nil
	}

	return nil, nil
}

func (t Table) buildInt(f field.Int, op Operator, vals stringValues) (field.Expr, error) {
	arg0, _ := vals.int0()
	switch op {
	case Eq:
		return f.Eq(arg0), nil
	case Neq:
		return f.Neq(arg0), nil
	case Gt:
		return f.Gt(arg0), nil
	case Gte:
		return f.Gte(arg0), nil
	case Lt:
		return f.Lt(arg0), nil
	case Lte:
		return f.Lte(arg0), nil
	case Like:
		return f.Like(arg0), nil
	case NotLike:
		return f.NotLike(arg0), nil
	case Between:
		arg1, _ := vals.intN(1)
		return f.Between(arg0, arg1), nil
	case NotBetween:
		arg1, _ := vals.intN(1)
		return f.NotBetween(arg0, arg1), nil
	case In:

	case NotIn:
	case Null:
	}

	return nil, nil
}

func (t Table) buildFloat64(f field.Float64, op Operator, vals stringValues) (field.Expr, error) {
	switch op {
	case Eq:
	case Neq:
	case Gt:
	case Gte:
	case Lt:
	case Lte:
	case Like:
	case NotLike:
	case Between:
	case NotBetween:
	case In:
	case NotIn:
	case Null:
	}

	return nil, nil
}

func (t Table) buildBool(f field.Bool, op Operator, vals stringValues) (field.Expr, error) {
	switch op {
	case Eq:
	case Neq:
	case Gt:
	case Gte:
	case Lt:
	case Lte:
	case Like:
	case NotLike:
	case Between:
	case NotBetween:
	case In:
	case NotIn:
	case Null:
	}

	return nil, nil
}

func (t Table) buildTime(f field.Time, op Operator, vals stringValues) (field.Expr, error) {
	switch op {
	case Eq:
	case Neq:
	case Gt:
	case Gte:
	case Lt:
	case Lte:
	case Like:
	case NotLike:
	case Between:
	case NotBetween:
	case In:
	case NotIn:
	case Null:
	}

	return nil, nil
}

type WhereBuild struct {
	Conds []gen.Condition
	Exprs []field.Expr
	Hits  []field.Expr // 命中的原始 field.Expr
}

type OrderBuild struct {
	Exprs []field.Expr
	Hits  []field.Expr
}
