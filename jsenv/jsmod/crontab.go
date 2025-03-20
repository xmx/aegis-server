package jsmod

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/xmx/aegis-server/jsenv/jsvm"
	"github.com/xmx/aegis-server/library/cronv3"
)

func NewCrontab(crond *cronv3.Crontab) jsvm.GlobalRegister {
	return &extCrontab{crond: crond}
}

type extCrontab struct {
	rt    jsvm.Runtime
	crond *cronv3.Crontab
}

func (ext *extCrontab) RegisterGlobal(rt jsvm.Runtime) error {
	ext.rt = rt
	fns := map[string]any{
		"addJob": ext.addJob,
	}

	return rt.Runtime().Set("crontab", fns)
}

func (ext *extCrontab) addJob(spec string, cmd func()) error {
	buf := make([]byte, 100)
	_, _ = rand.Read(buf)
	name := hex.EncodeToString(buf)
	ext.rt.AddFinalizer(&stopCron{name: name, crond: ext.crond})

	return ext.crond.AddJob(name, spec, ext.safeCall(cmd))
}

func (ext *extCrontab) safeCall(f func()) func() {
	return func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		f()
	}
}

type stopCron struct {
	name  string
	crond *cronv3.Crontab
}

func (s *stopCron) Finalize() error {
	s.crond.Remove(s.name)
	return nil
}
