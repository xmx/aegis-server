package dynsql

import (
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type Orders []*Order

type Order struct {
	ID     string     `json:"id"` // 唯一名字
	Table  string     `json:"table"`
	Column string     `json:"column"`
	Name   string     `json:"name"`
	Expr   field.Expr `json:"-"`
	db     *gorm.DB
}

func (o Order) Equals(fe field.Expr) bool {
	stmt := &gorm.Statement{DB: o.db}
	fe.Build(stmt)
	fname := stmt.SQL.String()
	stmt.SQL.Reset()
	o.Expr.Build(stmt)
	wname := stmt.SQL.String()

	return fname == wname
}

func (o Order) fullID() string {
	return o.Table + "." + o.Column
}
