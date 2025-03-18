package jsvm

import "github.com/grafana/sobek"

func New() *sobek.Runtime {
	vm := sobek.New()
	mapper := sobek.TagFieldNameMapper("json", true)
	vm.SetFieldNameMapper(mapper)
	vm.SetMaxCallStackSize(64)

	return vm
}

type Module interface {
	Register(root *Object) error
}

func RegisterModules(vm *sobek.Runtime, name string, mods []Module) error {
	obj := vm.NewObject()
	vm.Set(name, obj)

	root := &Object{vm: vm, obj: obj}

	for _, mod := range mods {
		if err := mod.Register(root); err != nil {
			return err
		}
	}

	return nil
}

type Object struct {
	vm  *sobek.Runtime
	obj *sobek.Object
}

func (o *Object) Sub(name string) *Object {
	sub := o.obj.Get(name)
	if obj, ok := sub.(*sobek.Object); ok {
		return &Object{vm: o.vm, obj: obj}
	}

	obj := o.vm.NewObject()
	_ = o.obj.Set(name, obj)

	return &Object{vm: o.vm, obj: obj}
}

func (o *Object) Set(name string, val any) {
	_ = o.obj.Set(name, val)
}

type GlobalRegister interface {
	RegisterGlobal(vm *sobek.Runtime) error
}

func RegisterGlobals(vm *sobek.Runtime, mods []GlobalRegister) error {
	for _, mod := range mods {
		if err := mod.RegisterGlobal(vm); err != nil {
			return err
		}
	}

	return nil
}
