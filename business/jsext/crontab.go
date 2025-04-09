package jsext

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/xmx/aegis-server/jsrun/jsvm"
	"github.com/xmx/aegis-server/library/cronv3"
)

func NewCrontab(crond *cronv3.Crontab) jsvm.GlobalRegister {
	return &extCrontab{crond: crond}
}

type extCrontab struct {
	rt    jsvm.Engineer
	crond *cronv3.Crontab
}

func (ext *extCrontab) RegisterGlobal(rt jsvm.Engineer) error {
	ext.rt = rt
	fns := map[string]any{
		"addJob": ext.addJob,
	}

	return rt.Runtime().Set("crontab", fns)
}

func (ext *extCrontab) addJob(spec string, cmd func()) error {
	buf := make([]byte, 30)
	_, _ = rand.Read(buf)
	name := hex.EncodeToString(buf)
	ext.rt.AddFinalizer(ext.remove(name))

	return ext.crond.AddJob(name, spec, ext.safeCall(cmd))
}

func (ext *extCrontab) safeCall(f func()) func() {
	return func() {
		defer func() {
			_ = recover()
		}()
		f()
	}
}

func (ext *extCrontab) remove(name string) func() error {
	return func() error {
		ext.crond.Remove(name)
		return nil
	}
}
