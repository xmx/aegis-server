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

// AddJob 添加定时任务（通过 cron 表达式），成功则会覆盖掉已存在的同名任务。
func (c *Crontab) AddJob(name string, spec string, cmd func()) error {
	exec := c.wrapFunc(name, cmd)
	attrs := []any{slog.String("name", name), slog.String("spec", spec)}

	c.mutex.Lock()
	oldID, exists := c.names[name]           // 先查询是否已存在任务
	newID, err := c.crond.AddJob(spec, exec) // 添加任务
	if err == nil {
		c.names[name] = newID // 添加成功，更新任务名字映射表
		if exists {
			c.crond.Remove(oldID) // 删除掉老的任务
		}
	}
	c.mutex.Unlock()

	if err != nil {
		return err
	} else if exists {
		c.log.Warn("修改定时任务", attrs...)
	} else {
		c.log.Info("新增定时任务", attrs...)
	}

	return err
}

// Schedule 添加定时任务，如果已存在的同名任务则会覆盖。
func (c *Crontab) Schedule(name string, sch cron.Schedule, job cron.Job) {
	exec := c.wrapFuncJob(name, job)

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
	if exists {
		c.log.Warn("删除定时任务", slog.String("name", name))
	}

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

// Cleanup 清理无效的任务，一般清理逻辑可以交给定时器定时清理。
//
// 对于常规大多数定时任务都是周而复始、一直能够执行的，但是有一些自定义 [cron.Schedule] 的任务
// 可能会返回无效时间，导致定时任务不再执行它，但是一直存在定时任务中。
// 还有一种场景就是 [NewSpecificTimes] 这种定时任务，只会在某几个时间点执行，完事后就无需再被
// 执行，最后一次定时任务执行完后，下次执行时间会返回空，定时器不再执行它，也需要清理。
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

func (c *Crontab) wrapFunc(name string, cmd func()) cron.FuncJob {
	return c.wrapFuncJob(name, cron.FuncJob(cmd))
}

func (c *Crontab) wrapFuncJob(name string, job cron.Job) cron.FuncJob {
	return func() {
		attrs := []any{slog.String("name", name)}
		c.log.Info("定时任务开始执行", attrs...)
		defer func() {
			if cause := recover(); cause != nil {
				attrs = append(attrs, slog.Any("cause", cause))
				c.log.Error("定时任务开始执行", attrs...)
			}
		}()
		job.Run()
		c.log.Info("定时任务执行结束", attrs...)
	}
}
