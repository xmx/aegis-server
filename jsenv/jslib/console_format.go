package jslib

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dop251/goja"
)

type consoleFormat struct{}

func (cf *consoleFormat) parse(f rune, val goja.Value, w *bytes.Buffer) bool {
	switch f {
	case 's':
		w.WriteString(val.String())
	case 'd':
		w.WriteString(val.ToNumber().String())
	case 'j':
		err := json.NewEncoder(w).Encode(val.Export())
		return err == nil
	case '%':
		w.WriteByte('%')
		return false
	default:
		w.WriteByte('%')
		w.WriteRune(f)
		return false
	}
	return true
}

func (cf *consoleFormat) formatTo(b *bytes.Buffer, f string, args ...goja.Value) {
	var pct bool
	var argNum int
	for _, chr := range f {
		if pct {
			if argNum < len(args) {
				if cf.parse(chr, args[argNum], b) {
					argNum++
				}
			} else {
				b.WriteByte('%')
				b.WriteRune(chr)
			}
			pct = false
		} else {
			if chr == '%' {
				pct = true
			} else {
				b.WriteRune(chr)
			}
		}
	}

	for _, arg := range args[argNum:] {
		b.WriteByte(' ')
		b.WriteString(arg.String())
	}
}

func (cf *consoleFormat) format(call goja.FunctionCall) []byte {
	buf := new(bytes.Buffer)
	//var format string
	//if arg := call.Argument(0); !goja.IsUndefined(arg) {
	//	format = arg.String()
	//}
	format := call.Argument(0).String()
	var args []goja.Value
	if len(call.Arguments) > 0 {
		args = call.Arguments[1:]
	}
	cf.formatTo(buf, format, args...)

	return buf.Bytes()
}

type formater interface {
	format(call goja.FunctionCall) []byte
}

type stdFormat struct{}

func (stdFormat) format(call goja.FunctionCall) []byte {
	//var format string
	//if arg := call.Argument(0); !goja.IsUndefined(arg) {
	//	format = arg.String()
	//}
	format := call.Argument(0).String()
	var vals []any
	if args := call.Arguments; len(args) != 0 {
		for _, arg := range args[1:] {
			val := arg.Export()
			vals = append(vals, val)
		}
	}
	buf := new(bytes.Buffer)
	_, _ = fmt.Fprintf(buf, format, vals...)

	return buf.Bytes()
}
