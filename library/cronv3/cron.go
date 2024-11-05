package cronv3

import (
	"context"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

func New(opts ...cron.Option) *Crontab {
	return &Crontab{
		crond: cron.New(opts...),
		names: make(map[string]cron.EntryID, 16),
	}
}

type Crontab struct {
	crond *cron.Cron
	mutex sync.Mutex
	names map[string]cron.EntryID
}

func (c *Crontab) Start() {
	c.crond.Start()
}

func (c *Crontab) Run() {
	c.crond.Run()
}

func (c *Crontab) Stop() context.Context {
	return c.crond.Stop()
}

func (c *Crontab) Location() *time.Location {
	return c.crond.Location()
}

func (c *Crontab) AddFunc(name, spec string, cmd func()) error {
	return c.AddJob(name, spec, cron.FuncJob(cmd))
}

func (c *Crontab) AddJob(name string, spec string, job cron.Job) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if id, ok := c.names[name]; ok {
		c.crond.Remove(id)
	}
	id, err := c.crond.AddJob(spec, job)
	if err == nil {
		c.names[name] = id
	}

	return err
}

func (c *Crontab) Schedule(name string, sch cron.Schedule, job cron.Job) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if id, ok := c.names[name]; ok {
		c.crond.Remove(id)
	}
	id := c.crond.Schedule(sch, job)
	c.names[name] = id
}

func (c *Crontab) Remove(name string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if id, ok := c.names[name]; ok {
		c.crond.Remove(id)
		delete(c.names, name)
	}
}
