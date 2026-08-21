package mongodsl

import "time"

type object interface {
	invoke(method string, args []any) (object, error)
	unwrap() any
}

type timeObject struct {
	value time.Time
	funcs functionMap
}

func newTimeObject(v time.Time) *timeObject {
	obj := &timeObject{value: v}
	obj.funcs = functionMap{
		"add": obj.add,
	}

	return obj
}

func (t *timeObject) invoke(method string, args []any) (object, error) {
	return t.funcs.call(method, args)
}

func (t *timeObject) unwrap() any {
	return t.value
}

func (t *timeObject) add(str string) (*timeObject, error) {
	du, err := time.ParseDuration(str)
	if err != nil {
		return nil, err
	}

	return newTimeObject(t.value.Add(du)), nil
}
