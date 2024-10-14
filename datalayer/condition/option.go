package condition

import "gorm.io/gen/field"

type ParserOptions struct {
	IgnoreOrder []field.Expr
	IgnoreWhere []field.Expr
}

func (o *ParserOptions) inIgnoreOrder(expr field.Expr) bool {
	if o == nil {
		return false
	}

	for _, f := range o.IgnoreOrder {
		if o.equals(expr, f) {
			return true
		}
	}

	return false
}

func (o *ParserOptions) inIgnoreWhere(expr field.Expr) bool {
	if o == nil {
		return false
	}

	for _, f := range o.IgnoreWhere {
		if o.equals(expr, f) {
			return true
		}
	}

	return false
}

func (o *ParserOptions) equals(f1, f2 field.Expr) bool {
	// FIXME 该方法只是比较 column 是否一样，并没有判断 table name。
	n1 := f1.ColumnName().String()
	n2 := f2.ColumnName().String()
	return n1 == n2
}
