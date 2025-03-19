package jsmod

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func NewHTTP() jsvm.GlobalRegister {
	return new(stdHTTP)
}

type stdHTTP struct{}

func (std *stdHTTP) RegisterGlobal(vm *goja.Runtime) error {
	fns := map[string]any{
		"post": std.post,
	}
	return vm.Set("http", fns)
}

func (std *stdHTTP) post(addr string, body any) {
	fmt.Println("post", addr)
	fmt.Println("post", body)
}
