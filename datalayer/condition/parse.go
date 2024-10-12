package condition

import (
	"reflect"

	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func ParseModel(db *gorm.DB, tbl any) (*Cond, error) {
	stmt := gorm.Statement{DB: db}
	if err := stmt.Parse(tbl); err != nil {
		return nil, err
	}

	sch := stmt.Schema
	table, fields := sch.Table, sch.Fields
	size := len(fields)
	conds := make(CondFields, 0, size)
	nameMap := make(map[string]*CondField, size)
	for _, f := range fields {
		expr := newField(table, f)
		name := f.DBName
		comment := f.Comment
		if comment == "" {
			comment = f.Name
		}
		cond := &CondField{name: name, comment: comment, expr: expr}
		conds = append(conds, cond)
	}
	ret := &Cond{fileds: conds, nameMap: nameMap}

	return ret, nil
}

// https://github.com/go-gorm/gen/blob/v0.3.26/internal/template/struct.go#L48
func newField(tbl string, f *schema.Field) field.Expr {
	name := f.DBName
	realType := getFieldRealType(f.FieldType)
	switch realType {
	case "string":
		return field.NewString(tbl, name)
	case "bytes":
		return field.NewBytes(tbl, name)
	case "int":
		return field.NewInt(tbl, name)
	case "int8":
		return field.NewInt8(tbl, name)
	case "int16":
		return field.NewInt16(tbl, name)
	case "int32":
		return field.NewInt32(tbl, name)
	case "int64":
		return field.NewInt64(tbl, name)
	case "uint":
		return field.NewUint(tbl, name)
	case "uint8":
		return field.NewUint8(tbl, name)
	case "uint16":
		return field.NewInt16(tbl, name)
	case "uin32":
		return field.NewUint32(tbl, name)
	case "uint64":
		return field.NewUint64(tbl, name)
	case "float32":
		return field.NewFloat32(tbl, name)
	case "float64":
		return field.NewFloat64(tbl, name)
	case "bool":
		return field.NewBool(tbl, name)
	case "time.Time":
		return field.NewTime(tbl, name)
	case "json.RawMessage", "[]byte":
		return field.NewBytes(tbl, name)
	case "serializer":
		return field.NewSerializer(tbl, name)
	default:
		return field.NewField(tbl, name)
	}
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
