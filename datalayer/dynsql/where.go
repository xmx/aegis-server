package dynsql

import (
	"strconv"

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
	opMaps    map[string]Operator
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

type WhereInputs []*WhereInput

type WhereInput struct {
	ID       string
	Operator Operator
	Values   []string
}

type stringValues []string

func (svs stringValues) string0() (string, bool) {
	return svs.stringN(0)
}

func (svs stringValues) stringN(i int) (string, bool) {
	if sz := len(svs); sz < i {
		return "", false
	}
	return svs[i], true
}

func (svs stringValues) bool0() (bool, bool) {
	if str, ok := svs.string0(); ok {
		val, err := strconv.ParseBool(str)
		return val, err == nil
	}

	return false, false
}

func (svs stringValues) int0() (int, bool) {
	return svs.intN(0)
}

func (svs stringValues) intN(i int) (int, bool) {
	if str, ok := svs.stringN(i); ok {
		val, err := strconv.Atoi(str)
		return val, err == nil
	}
	return 0, false
}

func (svs stringValues) ints() ([]int, error) {
	dats := make([]int, 0, 10)
	for _, sv := range svs {
		if n, err := strconv.Atoi(sv); err == nil {
			dats = append(dats, n)
		}
	}

	return dats, nil
}
