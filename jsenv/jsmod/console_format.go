package jsmod

import (
	"bytes"
	"encoding/json"

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
		err := json.NewEncoder(w).Encode(val)
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

func (cf *consoleFormat) format(b *bytes.Buffer, f string, args ...goja.Value) {
	pct := false
	argNum := 0
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

func (cf *consoleFormat) Format(call goja.FunctionCall) string {
	b := new(bytes.Buffer)

	var fmt string
	if arg := call.Argument(0); !goja.IsUndefined(arg) {
		fmt = arg.String()
	}

	var args []goja.Value
	if len(call.Arguments) > 0 {
		args = call.Arguments[1:]
	}
	cf.format(b, fmt, args...)

	return b.String()
}
