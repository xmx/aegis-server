package jsmod

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/dop251/goja"
	"github.com/xmx/aegis-server/jsrun/jsvm"
)

func NewConsole(stdout, stderr io.Writer) jsvm.ModuleRegister {
	return &writerConsole{stdout: stdout, stderr: stderr}
}

type writerConsole struct {
	stdout io.Writer
	stderr io.Writer
	vm     jsvm.Engineer
}

func (wc *writerConsole) RegisterModule(vm jsvm.Engineer) error {
	wc.vm = vm
	fields := map[string]any{
		"log":   wc.write,
		"error": wc.writeToStdout,
		"warn":  wc.write,
		"info":  wc.write,
		"debug": wc.write,
	}

	return vm.Runtime().Set("console", fields)
}

func (wc *writerConsole) write(call goja.FunctionCall) goja.Value {
	msg, err := wc.format(call)
	if err == nil {
		_, err = wc.stdout.Write(msg)
	}
	if err != nil {
		return wc.vm.Runtime().ToValue(err)
	}
	return goja.Undefined()
}

func (wc *writerConsole) writeToStdout(call goja.FunctionCall) goja.Value {
	msg, err := wc.format(call)
	if err == nil {
		_, err = wc.stderr.Write(msg)
	}
	if err != nil {
		return wc.vm.Runtime().ToValue(err)
	}
	return goja.Undefined()
}

func (wc *writerConsole) format(call goja.FunctionCall) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, arg := range call.Arguments {
		if err := wc.parse(buf, arg); err != nil {
			return nil, err
		}
	}
	buf.WriteByte('\n')

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
		_, _ = buf.ReadFrom(tmp)
	}

	return nil
}
