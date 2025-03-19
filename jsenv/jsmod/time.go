package jsmod

import (
	"time"

	"github.com/dop251/goja"
	"github.com/xmx/aegis-server/jsenv/jsvm"
)

func NewTime() jsvm.GlobalRegister {
	return new(stdTime)
}

type stdTime struct{}

func (*stdTime) RegisterGlobal(vm *goja.Runtime) error {
	fns := map[string]interface{}{
		"nanosecond ":   time.Nanosecond,
		"microsecond":   time.Microsecond,
		"millisecond":   time.Millisecond,
		"second":        time.Second,
		"minute":        time.Minute,
		"hour":          time.Hour,
		"january":       time.January,
		"february":      time.February,
		"march":         time.March,
		"april":         time.April,
		"may":           time.May,
		"june":          time.June,
		"july":          time.July,
		"august":        time.August,
		"september":     time.September,
		"october":       time.October,
		"november":      time.November,
		"december":      time.December,
		"sunday":        time.Sunday,
		"monday":        time.Monday,
		"tuesday":       time.Tuesday,
		"wednesday":     time.Wednesday,
		"thursday":      time.Thursday,
		"friday":        time.Friday,
		"saturday":      time.Saturday,
		"sleep":         time.Sleep,
		"local":         time.Local,
		"parseDuration": time.ParseDuration,
	}
	return vm.Set("time", fns)
}
