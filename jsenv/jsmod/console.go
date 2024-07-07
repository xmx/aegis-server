package jsmod

import (
	"io"
	"log"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func Console(w io.Writer) jsvm.Loader {
	if w == nil || w == io.Discard {
		return new(discardConsole)
	}

	l := log.New(w, "", log.LstdFlags)
	f := new(consoleFormat)

	return &printConsole{l: l, f: f}
}

type printConsole struct {
	l  *log.Logger
	f  *consoleFormat
	vm *goja.Runtime
}

func (p *printConsole) Global(vm *goja.Runtime) error {
	fields := map[string]any{
		"log":   p.print,
		"error": p.print,
		"warn":  p.print,
		"info":  p.print,
		"debug": p.print,
	}

	return vm.Set("console", fields)
}

func (p *printConsole) Require() (string, require.ModuleLoader) {
	return "", nil
}

func (p *printConsole) print(call goja.FunctionCall) goja.Value {
	msg := p.f.Format(call)
	p.l.Print(msg)
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
