package crontab

import (
	"context"
	"log/slog"
	"sync"

	"github.com/robfig/cron/v3"
)

type Cron struct {
	log   *slog.Logger
	ctab  *cron.Cron
	mutex sync.RWMutex
	tasks map[string]*scheduledTask
}

func New(log *slog.Logger, opts ...cron.Option) *Cron {
	return &Cron{
		log:   log,
		ctab:  cron.New(opts...),
		tasks: make(map[string]*scheduledTask, 8),
	}
}

func (c *Cron) Start() {
	c.ctab.Start()
}

func (c *Cron) Stop() {
	c.ctab.Stop()
}

func (c *Cron) Clean() []*EntryInfo {
	var ret []*EntryInfo
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for id, task := range c.tasks {
		entryID := task.entryID
		entry := c.ctab.Entry(entryID)
		if entry.Next.IsZero() {
			if info := c.remove(id); info != nil {
				ret = append(ret, info)
			}
		}
	}

	return ret
}

func (c *Cron) AddTask(task Tasker) error {
	info := task.Info()
	if info.ID == "" {
		info.ID = QualifiedID(task)
	}

	return c.addFunc(info, task.Call)
}

func (c *Cron) AddTasks(tasks []Tasker) error {
	for _, task := range tasks {
		if err := c.AddTask(task); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cron) AddJob(spec string, cmd func()) error {
	id := QualifiedID(cmd)
	info := TaskInfo{ID: id, CronSpec: spec}
	fn := wrapFunc(cmd)

	return c.addFunc(info, fn)
}

func (c *Cron) AddSchedule(sch cron.Schedule, cmd func()) error {
	id := QualifiedID(cmd)
	info := TaskInfo{ID: id, CronSched: sch}
	fn := wrapFunc(cmd)

	return c.addFunc(info, fn)
}

func (c *Cron) Remove(id string) *EntryInfo {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.remove(id)
}

func (c *Cron) Tasks() []*EntryInfo {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var ret []*EntryInfo
	for _, task := range c.tasks {
		entry := c.ctab.Entry(task.entryID)
		et := &EntryInfo{Entry: entry, Info: task.info}
		ret = append(ret, et)
	}

	return ret
}

func (c *Cron) remove(id string) *EntryInfo {
	if task := c.tasks[id]; task != nil {
		entry := c.ctab.Entry(task.entryID)
		c.ctab.Remove(task.entryID)
		delete(c.tasks, id)
		return &EntryInfo{
			Entry: entry,
			Info:  task.info,
		}
	}

	return nil
}

func (c *Cron) addFunc(info TaskInfo, exec func(context.Context) error) error {
	if info.ID == "" {
		info.ID = QualifiedID(exec)
	}

	id := info.ID
	sch := &scheduledTask{info: info, exec: exec, log: c.log}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 删除原来的任务（如果存在）。
	// 如果用户输入了错误的 cron 表达式，c.ctab.AddJob 添加失败，
	// 并且把老的 Task 也不会恢复。
	c.remove(id)

	if info.CronSched != nil {
		sch.entryID = c.ctab.Schedule(info.CronSched, sch)
	} else {
		entryID, err := c.ctab.AddJob(info.CronSpec, sch)
		if err != nil {
			return err
		}
		sch.entryID = entryID
	}
	c.tasks[id] = sch

	if info.Immediate { // 立即执行
		go sch.Run()
	}

	return nil
}
