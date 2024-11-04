package jsrt

import "github.com/dop251/goja"

type VM struct {
	vm *goja.Runtime
}

func New() *VM {
	return &VM{
		vm: goja.New(),
	}
}
