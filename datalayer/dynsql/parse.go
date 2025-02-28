package dynsql

import (
	"context"
	"reflect"

	"github.com/xmx/aegis-server/datalayer/query"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Options struct {
	Where    func(*Where) *Where
	Order    func(*Order) *Order
	NonWhere bool // 不解析 where
	NonOrder bool // 不解析 order
}

func Parse(qry *query.Query, mods []any, opt Options) (*Table, error) {
	if opt.Where == nil {
		opt.Where = func(w *Where) *Where { return w }
	}
	if opt.Order == nil {
		opt.Order = func(o *Order) *Order { return o }
	}

	queryCtx := qry.WithContext(context.Background())
	db := queryCtx.User.UnderlyingDB()

	var wheres Wheres
	var orders Orders
	whereMaps := make(map[string]*Where, 32)
	orderMaps := make(map[string]*Order, 32)

	for _, mod := range mods {
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(mod); err != nil {
			return nil, err
		}

		tableName, fields := stmt.Schema.Table, stmt.Schema.Fields
		for _, fld := range fields {
			expr, ops := parse(fld)
			if expr == nil {
				continue
			}

			column, comment := fld.DBName, fld.Comment
			if !opt.NonWhere {
				where := &Where{
					Table:     tableName,
					Column:    column,
					Type:      fld.GORMDataType,
					Name:      comment,
					Operators: ops,
					Expr:      expr,
					db:        db,
				}
				if w := opt.Where(where); w != nil && len(w.Operators) != 0 {
					shortID := w.Column
					if old := whereMaps[shortID]; old == nil {
						w.ID = shortID
						whereMaps[shortID] = w
					} else {
						delete(whereMaps, shortID)
						oldID := old.fullID()
						old.ID = oldID
						whereMaps[oldID] = w
						newID := w.fullID()
						w.ID = newID
						whereMaps[newID] = w
					}
					wheres = append(wheres, w)
				}
			}

			if opt.NonOrder {
				continue
			}
			order := &Order{
				Table:  tableName,
				Column: column,
				Name:   comment,
				Expr:   expr,
				db:     db,
			}
			if o := opt.Order(order); o != nil {
				shortID := o.Column
				if old := orderMaps[shortID]; old == nil {
					o.ID = shortID
					orderMaps[shortID] = o
				} else {
					delete(orderMaps, shortID)
					oldID := old.fullID()
					old.ID = oldID
					orderMaps[oldID] = old
					newID := o.fullID()
					o.ID = newID
					orderMaps[newID] = o
				}
				orders = append(orders, o)
			}
		}
	}

	tbl := &Table{
		wheres:    wheres,
		orders:    orders,
		whereMaps: whereMaps,
		orderMaps: orderMaps,
	}

	return tbl, nil
}

func parse(f *schema.Field) (field.Expr, Operators) {
	table, column := f.Schema.Table, f.DBName
	realType := getFieldRealType(f.FieldType)
	switch realType {
	case "string":
		return field.NewString(table, column),
			Operators{Eq, Neq, Gt, Gte, Lt, Lte, Like, NotLike, Between, NotBetween, In, NotIn, NotNull}
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64":
		return field.NewInt(table, column),
			Operators{Eq, Neq, Gt, Gte, Lt, Lte, Between, NotBetween, In, NotIn, NotNull}
	case "float32", "float64":
		return field.NewFloat64(table, column),
			Operators{Eq, Neq, Gt, Gte, Lt, Lte, Between, NotBetween, In, NotIn, NotNull}
	case "bool":
		return field.NewBool(table, column),
			Operators{Eq, Neq}
	case "time.Time":
		return field.NewTime(table, column),
			Operators{Eq, Neq, Gt, Gte, Lt, Lte, Between, NotBetween, In, NotIn, NotNull}
	default:
		return nil, nil
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

	str := f.String()
	if str == "time.Time" {
		return str
	}
	if str == "[]uint8" || str == "json.RawMessage" {
		return "bytes"
	}

	return f.Kind().String()
}
