package condition

import "gorm.io/gen/field"

type ParserOptions struct {
	IgnoreOrder []field.Expr
	IgnoreWhere []field.Expr
}
