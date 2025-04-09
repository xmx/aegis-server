package jsmod

import (
	"runtime"

	"github.com/xmx/aegis-server/jsrun/jsvm"
)

func NewRuntime() jsvm.GlobalRegister {
	return new(stdRuntime)
}

type stdRuntime struct{}

func (s *stdRuntime) RegisterGlobal(vm jsvm.Engineer) error {
	fns := map[string]any{
		"memStats": s.memStats,
		"goos":     runtime.GOOS,
		"goarch":   runtime.GOARCH,
	}
	return vm.Runtime().Set("runtime", fns)
}

func (s *stdRuntime) memStats() *runtime.MemStats {
	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)
	return stats
}
