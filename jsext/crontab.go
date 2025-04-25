package jsext

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"time"

	"github.com/dop251/goja"
	"github.com/xmx/aegis-server/library/cronv3"
	"github.com/xmx/jsos/jsvm"
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

func (ext *extCrontab) addJob(spec string, cmd func()) (*goja.Object, error) {
	buf := make([]byte, 16)
	nano := time.Now().UnixNano()
	binary.BigEndian.PutUint64(buf, uint64(nano))
	_, _ = rand.Read(buf[8:])
	name := hex.EncodeToString(buf)

	if err := ext.crond.AddJob(name, spec, cmd); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	task := &cronTask{
		name:   name,
		crond:  ext.crond,
		ctx:    ctx,
		cancel: cancel,
	}
	ext.eng.AddFinalizer(task.remove)
	taskObj := ext.eng.Runtime().NewObject()
	_ = taskObj.Set("name", task.name)
	_ = taskObj.Set("wait", task.wait)
	_ = taskObj.Set("remove", task.remove)

	return taskObj, nil
}

type cronTask struct {
	name   string
	crond  *cronv3.Crontab
	ctx    context.Context
	cancel context.CancelFunc
}

func (ct *cronTask) wait() {
	<-ct.ctx.Done()
}

func (ct *cronTask) remove() error {
	ct.crond.Remove(ct.name)
	ct.cancel()

	return nil
}
