package jsmod

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/grafana/sobek"
)

func NewConsole(w io.Writer) interface{ SetGlobal(vm *sobek.Runtime) error } {
	if w == nil || w == io.Discard {
		return new(discardConsole)
	}

	return &writerConsole{w: w}
}

type discardConsole struct{}

func (c *discardConsole) SetGlobal(vm *sobek.Runtime) error {
	fields := map[string]any{
		"log":   c.discord,
		"error": c.discord,
		"warn":  c.discord,
		"info":  c.discord,
		"debug": c.discord,
	}

	return vm.Set("console", fields)
}

func (*discardConsole) discord(sobek.FunctionCall) sobek.Value {
	return sobek.Undefined()
}

type writerConsole struct {
	w  io.Writer
	vm *sobek.Runtime
}

func (wc *writerConsole) SetGlobal(vm *sobek.Runtime) error {
	wc.vm = vm
	fields := map[string]any{
		"log":   wc.write,
		"error": wc.write,
		"warn":  wc.write,
		"info":  wc.write,
		"debug": wc.write,
	}

	return vm.Set("console", fields)
}

func (wc *writerConsole) write(call sobek.FunctionCall) sobek.Value {
	msg, err := wc.format(call)
	if err == nil {
		_, err = wc.w.Write(msg)
	}
	if err != nil {
		return wc.vm.ToValue(err)
	}
	return sobek.Undefined()
}

func (wc *writerConsole) format(call sobek.FunctionCall) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, arg := range call.Arguments {
		if err := wc.parse(buf, arg); err != nil {
			return nil, err
		}
	}
	buf.WriteByte('\n')

	return buf.Bytes(), nil
}

func (wc *writerConsole) parse(buf *bytes.Buffer, val sobek.Value) error {
	switch {
	case sobek.IsUndefined(val), sobek.IsNull(val):
		buf.WriteString(val.String())
		return nil
	}

	export := val.Export()
	switch v := export.(type) {
	case fmt.Stringer:
		buf.WriteString(v.String())
	case string:
		buf.WriteString(v)
	case int64:
		buf.WriteString(strconv.FormatInt(v, 10))
	case float64:
		buf.WriteString(strconv.FormatFloat(v, 'g', -1, 64))
	case bool:
		buf.WriteString(strconv.FormatBool(v))
	case []byte:
		str := base64.StdEncoding.EncodeToString(v)
		buf.WriteString(str)
	case func(sobek.FunctionCall) sobek.Value:
		buf.WriteString("<Function>")
	case sobek.ArrayBuffer:
		bs := v.Bytes()
		str := base64.StdEncoding.EncodeToString(bs)
		buf.WriteString(str)
	default:
		return wc.reflectParse(buf, v)
	}

	return nil
}

func (*writerConsole) reflectParse(buf *bytes.Buffer, v any) error {
	vof := reflect.ValueOf(v)
	switch vof.Kind() {
	case reflect.String:
		buf.WriteString(vof.String())
	case reflect.Int64:
		buf.WriteString(strconv.FormatInt(vof.Int(), 10))
	case reflect.Float64:
		buf.WriteString(strconv.FormatFloat(vof.Float(), 'g', -1, 64))
	case reflect.Bool:
		buf.WriteString(strconv.FormatBool(vof.Bool()))
	default:
		tmp := new(bytes.Buffer)
		if err := json.NewEncoder(tmp).Encode(v); err != nil || tmp.Len() == 0 {
			return err
		}
		_, _ = buf.ReadFrom(tmp)
	}

	return nil
}
