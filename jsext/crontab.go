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

func NewCrontab(crond *cronv3.Crontab) jsvm.ModuleLoader {
	return &crontabLoader{crond: crond}
}

type crontabLoader struct {
	crond *cronv3.Crontab
}

func (ld crontabLoader) LoadModule(eng jsvm.Engineer) error {
	cs := crontabSandbox{eng: eng, crond: ld.crond}
	vals := map[string]any{
		"addJob": cs.addJob,
	}
	eng.RegisterModule("crontab", vals, true)

	return nil
}

type crontabSandbox struct {
	eng   jsvm.Engineer
	crond *cronv3.Crontab
}

func (cs crontabSandbox) addJob(spec string, cmd func()) (*goja.Object, error) {
	buf := make([]byte, 16)
	nano := time.Now().UnixNano()
	binary.BigEndian.PutUint64(buf, uint64(nano))
	_, _ = rand.Read(buf[8:])
	name := hex.EncodeToString(buf)

	if _, err := cs.crond.AddJob(name, spec, cmd); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	task := &cronTask{
		name:   name,
		crond:  cs.crond,
		ctx:    ctx,
		cancel: cancel,
	}
	cs.eng.AddFinalizer(task.remove)
	taskObj := cs.eng.Runtime().NewObject()
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
