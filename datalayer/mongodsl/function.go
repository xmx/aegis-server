package mongodsl

import (
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type functionMap map[string]any

func (fm functionMap) invoke(name string, args []any) (any, error) {
	fv, exists := fm[name]
	if !exists {
		return nil, newDSLErrorf("方法 %s 不存在", name)
	}

	fval := reflect.ValueOf(fv)
	ftyp := fval.Type()
	numIn := ftyp.NumIn()
	if num := len(args); num != numIn {
		return nil, newDSLErrorf("方法 %s 期望得到 %d 个参数，但实际得到 %d 个", name, numIn, num)
	}

	inargs := make([]reflect.Value, numIn)
	for i, arg := range args {
		expect := ftyp.In(i)
		input := reflect.ValueOf(arg)

		if input.IsValid() && input.Type() != expect {
			if input.Type().ConvertibleTo(expect) {
				input = input.Convert(expect)
			} else {
				return nil, newDSLErrorf("方法 %s 的第 %d 个入参类型不符合期望", name, i)
			}
		}

		inargs[i] = input
	}

	results := fval.Call(inargs)
	numout := len(results)
	if numout == 0 {
		return nil, newDSLErrorf("方法 %s 必须要有返回值", name)
	}

	if rerr, yes := results[numout-1].Interface().(error); yes && rerr != nil {
		return nil, rerr
	}

	return results[0].Interface(), nil
}

func (fm functionMap) call(method string, args []any) (object, error) {
	ret, err := fm.invoke(method, args)
	if err != nil {
		return nil, err
	}

	obj, ok := ret.(object)
	if !ok {
		return nil, newDSLErrorf("方法 %s 的返回值类型不合法", method)
	}

	return obj, nil
}

func builtinBSONExpressionFunctions() functionMap {
	return functionMap{
		"now": func() *timeObject {
			return newTimeObject(time.Now())
		},
		"oid": func(s string) (bson.ObjectID, error) {
			return bson.ObjectIDFromHex(s)
		},
		"exists": func(s string) bson.M {
			return bson.M{s: bson.M{"$exists": true}}
		},
		"date": func(s string) (*timeObject, error) {
			var t time.Time
			var err error

			layouts := []string{time.RFC3339, time.RFC3339Nano, time.DateTime, time.DateOnly}
			for _, layout := range layouts {
				if t, err = time.Parse(layout, s); err == nil {
					break
				}
			}
			if err != nil {
				return nil, err
			}

			return newTimeObject(t), nil
		},
	}
}
