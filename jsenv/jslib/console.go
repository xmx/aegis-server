package jslib

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func Console(w io.Writer) jsvm.Loader {
	if w == nil || w == io.Discard {
		return new(discardConsole)
	}

	return &writerConsole{
		w: w,
	}
}

type discardConsole struct{}

func (c *discardConsole) Global(vm *goja.Runtime) error {
	fields := map[string]any{
		"log":   c.discord,
		"error": c.discord,
		"warn":  c.discord,
		"info":  c.discord,
		"debug": c.discord,
	}

	return vm.Set("console", fields)
}

func (c *discardConsole) Require() (string, require.ModuleLoader) {
	return "", nil
}

func (*discardConsole) discord(goja.FunctionCall) goja.Value {
	return goja.Undefined()
}

type writerConsole struct {
	w  io.Writer
	vm *goja.Runtime
}

func (wc *writerConsole) Global(vm *goja.Runtime) error {
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

func (wc *writerConsole) Require() (string, require.ModuleLoader) {
	return "", nil
}

func (wc *writerConsole) write(call goja.FunctionCall) goja.Value {
	msg, err := wc.format(call)
	if err == nil {
		_, err = wc.w.Write(msg)
	}
	if err != nil {
		return wc.vm.ToValue(err)
	}
	return goja.Undefined()
}

func (wc *writerConsole) format(call goja.FunctionCall) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, val := range call.Arguments {
		if err := wc.parse(buf, val); err != nil {
			return nil, err
		}
		buf.WriteByte(' ')
	}

	return buf.Bytes(), nil
}

func (wc *writerConsole) parse(buf *bytes.Buffer, val goja.Value) error {
	switch {
	case goja.IsUndefined(val), goja.IsNull(val):
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
	case func(goja.FunctionCall) goja.Value:
		buf.WriteString("<Function>")
	case goja.ArrayBuffer:
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
		// 标准库 JSON 序列化后会在最后面加换行符，这样输出到前端会有个换行，
		// 下面的操作就是去除 JSON 字符串后面的换行符。当然不处理也无伤大雅。
		_, _ = buf.ReadFrom(io.LimitReader(tmp, int64(tmp.Len()-1)))
	}

	return nil
}
