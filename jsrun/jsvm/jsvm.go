package jsvm

import (
	"sync"

	"github.com/dop251/goja"
)

func New(mods []ModuleRegister) (Engineer, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(newFieldNameMapper("json"))

	rqu := &require{
		modules: make(map[string]goja.Value, 16),
		sources: make(map[string]goja.Value, 16),
	}
	eng := &jsEngine{
		vm:      vm,
		require: rqu,
	}
	rqu.eng = eng
	_ = vm.Set("require", rqu.load)
	if err := RegisterModules(eng, mods); err != nil {
		return nil, err
	}

	return eng, nil
}

type jsEngine struct {
	vm      *goja.Runtime
	require *require
	mutex   sync.Mutex
	finals  []func() error
}

func (jse *jsEngine) Runtime() *goja.Runtime {
	return jse.vm
}

func (jse *jsEngine) RunString(code string) (goja.Value, error) {
	cjs, err := Transform(code, map[string]any{"plugins": []string{"transform-modules-commonjs"}})
	if err != nil {
		return nil, err
	}
	return jse.vm.RunString(cjs)
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
