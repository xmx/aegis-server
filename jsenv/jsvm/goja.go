package jsvm

import (
	"crypto/rand"
	"encoding/hex"
	"sync"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

func New() *goja.Runtime {
	vm := goja.New()
	mapper := goja.TagFieldNameMapper("json", true)
	vm.SetFieldNameMapper(mapper)
	vm.SetMaxCallStackSize(64)
	detectOrRegister(vm)

	return vm
}

type Loader interface {
	Global(vm *goja.Runtime) error
	Require() (string, require.ModuleLoader)
}

func Register(vm *goja.Runtime, loaders []Loader) error {
	registry := detectOrRegister(vm)
	for _, l := range loaders {
		if l == nil {
			continue
		}
		if err := l.Global(vm); err != nil {
			return err
		}
		name, requ := l.Require()
		if name == "" || requ == nil {
			continue
		}

		pxy := proxyRegistry(requ)
		registry.register(name, pxy.Call)
	}

	return nil
}

func detectOrRegister(vm *goja.Runtime) *hideRegistry {
	registryID := onceRandomRegistryID()
	if val := vm.Get(registryID); val != nil {
		if registry, ok := val.Export().(*hideRegistry); ok && registry != nil {
			return registry
		}
	}

	noop := func(string) ([]byte, error) { return nil, require.ModuleFileDoesNotExistError }
	registry := require.NewRegistryWithLoader(noop)
	registry.Enable(vm)
	hide := &hideRegistry{r: registry}
	_ = vm.Set(registryID, hide)

	return hide
}

// onceRandomRegistryID 程序每次随机生成一个 require 模块 ID。
// 防止脚本层发现调用或恶意篡改。
var onceRandomRegistryID = sync.OnceValue[string](func() string {
	buf := make([]byte, 20)
	_, _ = rand.Read(buf)
	registryID := ".REGISTRY-KEY-" + hex.EncodeToString(buf)

	return registryID
})

type hideRegistry struct {
	r *require.Registry
}

func (h *hideRegistry) register(name string, loader require.ModuleLoader) {
	h.r.RegisterNativeModule(name, loader)
}

type proxyRegistry func(*goja.Runtime, *goja.Object)

func (pr proxyRegistry) Call(vm *goja.Runtime, mod *goja.Object) {
	obj := mod.Get("exports").(*goja.Object)
	pr(vm, obj)
}
