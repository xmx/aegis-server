package condition

import (
	"reflect"

	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func ParseModel(db *gorm.DB, tbl any, opts *ParserOptions) (*Cond, error) {
	stmt := gorm.Statement{DB: db}
	if err := stmt.Parse(tbl); err != nil {
		return nil, err
	}

	sch := stmt.Schema
	table, fields := sch.Table, sch.Fields
	orders := make([]*OrderField, 0, 10)
	wheres := make([]*WhereField, 0, 10)
	ordersNameMap := make(map[string]*OrderField, 8)
	wheresNameMap := make(map[string]*WhereField, 8)
	for _, f := range fields {
		expr, ops := newField(table, f)
		name, comment := f.DBName, f.Comment
		if comment == "" {
			comment = f.Name
		}

		if exp := parseOrderField(expr, opts); exp != nil {
			cond := &OrderField{name: name, comment: comment, expr: exp}
			orders = append(orders, cond)
			ordersNameMap[name] = cond
		}
		if exp := parseWhereField(expr, opts); exp != nil {
			opMap := ops.NameMap()
			cond := &WhereField{name: name, comment: comment, expr: exp, operators: ops, opMap: opMap}
			wheres = append(wheres, cond)
			wheresNameMap[name] = cond
		}
	}
	ret := &Cond{
		orders:        orders,
		ordersNameMap: ordersNameMap,
		wheres:        wheres,
		wheresNameMap: wheresNameMap,
	}

	return ret, nil
}

func parseOrderField(f field.Expr, opts *ParserOptions) field.OrderExpr {
	if expr, ok := f.(field.OrderExpr); ok && !opts.inIgnoreOrder(f) {
		return expr
	}
	return nil
}

func parseWhereField(f field.Expr, opts *ParserOptions) field.Expr {
	if opts.inIgnoreWhere(f) {
		return nil
	}
	return f
}

// https://github.com/go-gorm/gen/blob/v0.3.26/internal/template/struct.go#L48
func newField(tbl string, f *schema.Field) (field.Expr, operators) {
	name := f.DBName
	var expr field.Expr
	realType := getFieldRealType(f.FieldType)
	ops := typeAllowedOperator(realType)
	switch realType {
	case "string":
		expr = field.NewString(tbl, name)
	case "int":
		expr = field.NewInt(tbl, name)
	case "int8":
		expr = field.NewInt8(tbl, name)
	case "int16":
		expr = field.NewInt16(tbl, name)
	case "int32":
		expr = field.NewInt32(tbl, name)
	case "int64":
		expr = field.NewInt64(tbl, name)
	case "uint":
		expr = field.NewUint(tbl, name)
	case "uint8":
		expr = field.NewUint8(tbl, name)
	case "uint16":
		expr = field.NewUint16(tbl, name)
	case "uin32":
		expr = field.NewUint32(tbl, name)
	case "uint64":
		expr = field.NewUint64(tbl, name)
	case "float32":
		expr = field.NewFloat32(tbl, name)
	case "float64":
		expr = field.NewFloat64(tbl, name)
	case "bool":
		expr = field.NewBool(tbl, name)
	case "time.Time":
		expr = field.NewTime(tbl, name)
	case "bytes", "[]byte", "json.RawMessage":
		expr = field.NewBytes(tbl, name)
	case "serializer":
		expr = field.NewSerializer(tbl, name)
	default:
		expr = field.NewField(tbl, name)
	}

	return expr, ops
}

func typeAllowedOperator(realType string) operators {
	switch realType {
	case "string", "bytes", "[]byte", "json.RawMessage":
		return operators{Eq, Neq, Gt, Gte, Lt, Lte, Like, NotLike, Regex, NotRegex, Between, NotBetween, In, NotIn}
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uin32", "uint64",
		"float32", "float64":
		return operators{Eq, Neq, Gt, Gte, Lt, Lte, Like, NotLike, Between, NotBetween, In, NotIn}
	case "bool":
		return operators{Eq, Neq}
	case "time.Time":
		return operators{Eq, Neq, Gt, Gte, Lt, Lte, Between, NotBetween, In, NotIn}
	case "serializer":
		return operators{Eq, Neq, Gt, Gte, Lt, Lte, Like, Regex, In}
	}
	return operators{Eq, Neq, Gt, Gte, Lt, Lte, Like, In}
}

// getFieldRealType  get basic type of field.
// https://github.com/go-gorm/gen/blob/v0.3.26/internal/generate/query.go#L75-L96
func getFieldRealType(f reflect.Type) string {
	serializerInterface := reflect.TypeOf((*schema.SerializerInterface)(nil)).Elem()
	if f.Implements(serializerInterface) || reflect.New(f).Type().Implements(serializerInterface) {
		return "serializer"
	}
	scanValuer := reflect.TypeOf((*field.ScanValuer)(nil)).Elem()
	if f.Implements(scanValuer) || reflect.New(f).Type().Implements(scanValuer) {
		return "field"
	}

	if f.Kind() == reflect.Ptr {
		f = f.Elem()
	}
	if f.String() == "time.Time" {
		return "time.Time"
	}
	if f.String() == "[]uint8" || f.String() == "json.RawMessage" {
		return "bytes"
	}
	return f.Kind().String()
}
