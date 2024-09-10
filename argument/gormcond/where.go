package gormcond

import (
	"strings"

	"gorm.io/gen/field"
)

type Where struct{}

//  =          eq       1 string int bool time
// !=          ne       1 string int bool time
//  >          gt       1 string int      time
//  <          lt       1 string int      time
// >=          gte      1 string int      time
// <=          lte      1 string int      time
// IN          in       n string int      time
// NOT IN      notin    n string int      time
// LIKE        like     1 string
// NOT LIKE    notlike  1 string
// BETWEEN     btw      2 string int      time
// NOT BETWEEN notbtw   2 string int      time
// REGEX       regex    1 string
// NOT REGEX   notregex 1 string

type Hi interface {
	field.Bool | field.String | field.Int | field.Time | field.Float64
}

func Ass(f field.Expr) {
	// string bool number
	switch f.(type) {
	case field.Bool:
	case field.String:
	case field.Int:
	}
}

func Parse(str string) {
	sn := strings.SplitN(str, ":", 3)
}
