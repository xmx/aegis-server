package jsext

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/xmx/aegis-server/jsrun/jsvm"
	"github.com/xmx/aegis-server/library/cronv3"
)

func NewCrontab(crond *cronv3.Crontab) jsvm.ModuleRegister {
	return &extCrontab{crond: crond}
}

type extCrontab struct {
	eng   jsvm.Engineer
	crond *cronv3.Crontab
}

func (ext *extCrontab) RegisterModule(eng jsvm.Engineer) error {
	ext.eng = eng
	vals := map[string]any{
		"addJob": ext.addJob,
	}
	eng.RegisterModule("crontab", vals, true)

	return nil
}

func (ext *extCrontab) addJob(spec string, cmd func()) error {
	buf := make([]byte, 30)
	_, _ = rand.Read(buf)
	name := hex.EncodeToString(buf)
	ext.eng.AddFinalizer(ext.remove(name))

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
