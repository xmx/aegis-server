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
					opMaps := make(map[string]Operator, len(w.Operators))
					for _, o := range w.Operators {
						id, _ := o.Info()
						opMaps[id] = o
					}
					w.db, w.opMaps = db, opMaps

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
			orderExpr, ok := expr.(field.OrderExpr)
			if !ok {
				continue
			}
			order := &Order{
				Table:  tableName,
				Column: column,
				Name:   comment,
				Expr:   orderExpr,
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
	var expr field.Expr
	var ops Operators

	table, column := f.Schema.Table, f.DBName
	realType := getFieldRealType(f.FieldType)
	switch realType {
	case "string":
		expr = field.NewString(table, column)
		ops = Operators{Eq, Neq, Gt, Gte, Lt, Lte, Like, NotLike, Between, NotBetween, In, NotIn}
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64":
		expr = field.NewInt(table, column)
		ops = Operators{Eq, Neq, Gt, Gte, Lt, Lte, Between, NotBetween, In, NotIn}
	case "float32", "float64":
		expr = field.NewFloat64(table, column)
		ops = Operators{Eq, Neq, Gt, Gte, Lt, Lte, Between, NotBetween, In, NotIn}
	case "bool":
		expr = field.NewBool(table, column)
		ops = Operators{Eq, Neq}
	case "time.Time":
		expr = field.NewTime(table, column)
		ops = Operators{Eq, Neq, Gt, Gte, Lt, Lte, Between, NotBetween, In, NotIn}
	default:
		return nil, nil
	}
	if !f.NotNull {
		ops = append(ops, Null)
	}

	return expr, ops
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
