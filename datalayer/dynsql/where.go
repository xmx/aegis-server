package dynsql

import (
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Wheres []*Where

type Where struct {
	ID        string          `json:"id"`        // 唯一名字
	Table     string          `json:"table"`     // 表名
	Column    string          `json:"column"`    // 列名
	Type      schema.DataType `json:"type"`      // 类型
	Name      string          `json:"name"`      // 注释
	Operators Operators       `json:"operators"` // 允许的操作符
	Enums     Enums           `json:"enums"`     // 枚举值
	Expr      field.Expr      `json:"-"`         // 列
	db        *gorm.DB
}

func (w Where) Equals(fe field.Expr) bool {
	stmt := &gorm.Statement{DB: w.db}
	fe.Build(stmt)
	fname := stmt.SQL.String()
	stmt.SQL.Reset()
	w.Expr.Build(stmt)
	wname := stmt.SQL.String()

	return fname == wname
}

func (w Where) fullID() string {
	return w.Table + "." + w.Column
}
