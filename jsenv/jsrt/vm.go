package jsrt

import "github.com/grafana/sobek"

func New() *sobek.Runtime {
	vm := sobek.New()
	mapper := sobek.TagFieldNameMapper("json", true)
	vm.SetFieldNameMapper(mapper)
	vm.SetMaxCallStackSize(64)

	return vm
}
