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

	return &writerConsole{
		w: w,
		f: new(normalFormat),
	}
}

type writerConsole struct {
	w  io.Writer
	f  formatter
	vm *goja.Runtime
}

func (p *writerConsole) Global(vm *goja.Runtime) error {
	p.vm = vm
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
	msg, err := p.f.format(call)
<<<<<<< HEAD
	if err == nil {
		_, err = p.w.Write(msg)
	}
	if err != nil {
		return p.vm.ToValue(err)
	}
=======
	if err != nil {
		return p.vm.ToValue(err)
	}
	if _, err = p.w.Write(msg); err != nil {
		return p.vm.ToValue(err)
	}
>>>>>>> 9f53d52 (go get -u)

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
