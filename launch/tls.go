package launch

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/datalayer/model"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func saveTLS(ctx context.Context, svc service.ConfigCertificate) {
	key, _ := os.ReadFile("resources/temp/lo.zzu.wiki.key")
	pem, _ := os.ReadFile("resources/temp/lo.zzu.wiki.pem")

	req := &request.ConfigCertificateCreate{
		PublicKey:  string(pem),
		PrivateKey: string(key),
		Enabled:    true,
	}
	svc.Create(ctx, req)
}

func parseModel(gdb *gorm.DB) {
	stmt := gorm.Statement{DB: gdb}
	_ = stmt.Parse(model.GridChunk{})
	sch := stmt.Schema
	for _, f := range sch.Fields {
		fmt.Println(f.Comment)
	}
}

// https://github.com/go-gorm/gen/blob/v0.3.26/internal/model/base.go#L193-L220
func newFields(sch *schema.Schema) []field.Expr {
	tbl := sch.Table
	exprs := make([]field.Expr, 0, len(sch.Fields))
	for _, f := range sch.Fields {
		expr := newField(tbl, f)
		exprs = append(exprs, expr)
	}
	return exprs
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
