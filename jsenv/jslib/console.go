package jslib

import (
	"io"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func Console(w io.Writer) jsvm.Loader {
	if w == nil || w == io.Discard {
		return new(discardConsole)
	}

	return &writerConsole{w: w}
}

type writerConsole struct {
	w io.Writer
	f consoleFormat
}

func (p *writerConsole) Global(vm *goja.Runtime) error {
	fields := map[string]any{
		"log":   p.print,
		"error": p.print,
		"warn":  p.print,
		"info":  p.print,
		"debug": p.print,
	}

	return vm.Set("console", fields)
}

func (p *writerConsole) Require() (string, require.ModuleLoader) {
	return "", nil
}

func (p *writerConsole) print(call goja.FunctionCall) goja.Value {
	msg := p.f.Format(call)
	_, _ = p.w.Write(msg)
	return goja.Undefined()
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

func (*discardConsole) discord(_ goja.FunctionCall) goja.Value {
	return goja.Undefined()
}
