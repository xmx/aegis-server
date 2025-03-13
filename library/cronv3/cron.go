package cronv3

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

func New(log *slog.Logger, opts ...cron.Option) *Crontab {
	return &Crontab{
		log:   log,
		crond: cron.New(opts...),
		names: make(map[string]cron.EntryID, 16),
	}
}

type Crontab struct {
	log   *slog.Logger
	crond *cron.Cron
	mutex sync.Mutex
	names map[string]cron.EntryID
}

func (c *Crontab) Start() {
	c.log.Warn("异步运行 crontab")
	c.crond.Start()
}

func (c *Crontab) Run() {
	c.log.Warn("同步运行 crontab")
	c.crond.Run()
}

func (c *Crontab) Stop() context.Context {
	c.log.Warn("停止 crontab")
	return c.crond.Stop()
}

func (c *Crontab) Location() *time.Location {
	return c.crond.Location()
}

func (c *Crontab) AddJob(name string, spec string, cmd func()) error {
	exec := c.withLogFunc(name, cmd)
	attrs := []any{slog.String("name", name), slog.String("spec", spec)}

	c.mutex.Lock()
	exists := c.remove(name)
	id, err := c.crond.AddJob(spec, exec)
	if err == nil {
		c.names[name] = id
	}
	c.mutex.Unlock()
	if err != nil {
		if exists {
			c.log.Warn("修改定时任务失败导致原来的任务被删除", attrs...)
		}
		return err
	}

	if exists {
		c.log.Warn("修改定时任务", attrs...)
	} else {
		c.log.Info("新增定时任务", attrs...)
	}

	return err
}

func (c *Crontab) Schedule(name string, sch cron.Schedule, job cron.Job) {
	exec := c.withLogJob(name, job)

	c.mutex.Lock()
	exists := c.remove(name)
	id := c.crond.Schedule(sch, exec)
	c.names[name] = id
	c.mutex.Unlock()

	attrs := []any{slog.String("name", name)}
	if exists {
		c.log.Warn("修改定时任务", attrs...)
	} else {
		c.log.Info("新增定时任务", attrs...)
	}
}

func (c *Crontab) Remove(name string) bool {
	c.mutex.Lock()
	exists := c.remove(name)
	c.mutex.Unlock()

	return exists
}

func (c *Crontab) remove(name string) bool {
	if id, ok := c.names[name]; ok {
		c.crond.Remove(id)
		delete(c.names, name)
		return true
	}

	return false
}

func (c *Crontab) Cleanup() {
	for _, ent := range c.crond.Entries() {
		if !ent.Next.IsZero() {
			continue
		}
		if name := c.lookupName(ent.ID); name != "" {
			c.Remove(name)
		}
	}
}

func (c *Crontab) lookupName(id cron.EntryID) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for name, eid := range c.names {
		if id == eid {
			return name
		}
	}

	return ""
}

func (c *Crontab) withLogFunc(name string, cmd func()) cron.FuncJob {
	return c.withLogJob(name, cron.FuncJob(cmd))
}

func (c *Crontab) withLogJob(name string, job cron.Job) cron.FuncJob {
	return func() {
		c.log.Info("定时任务开始执行", slog.String("name", name))
		job.Run()
		c.log.Info("定时任务执行结束", slog.String("name", name))
	}
}
