package jsmod

import (
	"runtime"

	"github.com/dop251/goja"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func NewRuntime() jsvm.GlobalRegister {
	return new(stdRuntime)
}

type stdRuntime struct{}

func (s *stdRuntime) RegisterGlobal(vm *goja.Runtime) error {
	fns := map[string]any{
		"memStats": s.memStats,
	}
	return vm.Set("runtime", fns)
}

func (s *stdRuntime) memStats() *runtime.MemStats {
	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)
	return stats
}
