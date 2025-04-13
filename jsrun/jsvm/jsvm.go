package jsvm

import (
	"archive/zip"
	"io"
	"sync"

	"github.com/dop251/goja"
)

func New(mods []ModuleRegister) (Engineer, error) {
	prog, err := onceCompileBabel()
	if err != nil {
		return nil, err
	}

	vm := goja.New()
	// babel need
	logFunc := func(goja.FunctionCall) goja.Value { return nil }
	_ = vm.Set("console", map[string]func(goja.FunctionCall) goja.Value{
		"log":   logFunc,
		"error": logFunc,
		"warn":  logFunc,
	})

	if _, err = vm.RunProgram(prog); err != nil {
		return nil, err
	}
	var transformFunc goja.Callable
	babel := vm.Get("Babel")
	if err = vm.ExportTo(babel.ToObject(vm).Get("transform"), &transformFunc); err != nil {
		return nil, err
	}

	transform := func(code string, opts map[string]any) (string, error) {
		if value, exx := transformFunc(babel, vm.ToValue(code), vm.ToValue(opts)); exx != nil {
			return "", exx
		} else {
			return value.ToObject(vm).Get("code").String(), nil
		}
	}
	eng := &jsEngine{
		vm:        vm,
		transform: transform,
	}
	rqu := &require{
		eng:     eng,
		modules: make(map[string]goja.Value, 16),
		sources: make(map[string]goja.Value, 16),
	}
	eng.require = rqu
	_ = vm.Set("require", rqu.load)

	if err = RegisterModules(eng, mods); err != nil {
		return nil, err
	}

	return eng, nil
}

type jsEngine struct {
	vm *goja.Runtime
	// transform Babel.transform()
	transform func(code string, opts map[string]any) (string, error)

	require *require

	mutex  sync.Mutex
	finals []func() error
}

func (jse *jsEngine) RunZip(zrd *zip.ReadCloser) (goja.Value, error) {
	file, err := zrd.Open("main.js")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	jse.require.source = zrd

	return jse.RunString(string(data))
}

func (jse *jsEngine) Runtime() *goja.Runtime {
	return jse.vm
}

func (jse *jsEngine) RunString(code string) (goja.Value, error) {
	commonJS, err := jse.transform(code, map[string]any{"plugins": []string{"transform-modules-commonjs"}})
	if err != nil {
		return nil, err
	}

	return jse.vm.RunString(commonJS)
}

func (jse *jsEngine) RunProgram(pgm *goja.Program) (goja.Value, error) {
	return jse.vm.RunProgram(pgm)
}

func (jse *jsEngine) RegisterModule(name string, module any, override bool) bool {
	return jse.require.register(name, module, override)
}

func (jse *jsEngine) AddFinalizer(finals ...func() error) {
	jse.mutex.Lock()
	defer jse.mutex.Unlock()

	for _, final := range finals {
		if final != nil {
			jse.finals = append(jse.finals, final)
		}
	}
}

func (jse *jsEngine) Interrupt(v any) {
	jse.vm.Interrupt(v)

	jse.mutex.Lock()
	defer jse.mutex.Unlock()

	for _, final := range jse.finals {
		_ = final()
	}
	jse.finals = nil
}

func (jse *jsEngine) ClearInterrupt() {
	jse.vm.ClearInterrupt()
}
