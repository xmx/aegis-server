package gopool

import (
	"context"
	"sync/atomic"
)

type Executor interface {
	Exec(ctx context.Context, task func(ctx context.Context))
	Await(parent context.Context, tasks []func(ctx context.Context)) context.Context
}

func Limit(work, queue int) Executor {
	if work < 1 {
		work = 1
	}
	if queue < 0 {
		queue = 0
	}

	return &limitExecutor{
		sem:   make(chan struct{}, work),
		queue: make(chan *taskFunc, queue),
	}
}

type limitExecutor struct {
	sem   chan struct{}
	queue chan *taskFunc
}

func (le *limitExecutor) Exec(ctx context.Context, task func(ctx context.Context)) {
	if task == nil {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}

	tf := &taskFunc{ctx: ctx, task: task}
	select {
	case le.sem <- struct{}{}:
		go tf.call()
		return
	default:
	}

	select {
	case <-ctx.Done():
	case le.queue <- tf:
	case le.sem <- struct{}{}:
		go tf.call()
	}
}

func (le *limitExecutor) Await(parent context.Context, tasks []func(ctx context.Context)) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	size := len(tasks)
	ctx, cancel := context.WithCancel(parent)
	if size == 0 {
		cancel()
		return ctx
	}

	count := &awaitCount{cancel: cancel}
	for _, task := range tasks {
		if task == nil {
			continue
		}

		count.incr()
		fn := le.awaitCountFunc(count, task)
		le.Exec(ctx, fn)
	}

	return ctx
}

func (le *limitExecutor) exec(tf *taskFunc) {
	defer func() { <-le.sem }()

	tf.call()

	var none bool
	for !none {
		select {
		case tfn := <-le.queue:
			tfn.call()
		default:
			none = true
		}
	}
}

func (le *limitExecutor) awaitCountFunc(cnt *awaitCount, task func(context.Context)) func(context.Context) {
	return func(ctx context.Context) {
		defer func() { cnt.decr() }()
		task(ctx)
	}
}

type taskFunc struct {
	ctx  context.Context
	task func(ctx context.Context)
}

func (tf *taskFunc) call() {
	defer func() { recover() }()
	tf.task(tf.ctx)
}

type awaitCount struct {
	count  atomic.Int64
	cancel context.CancelFunc
}

func (ac *awaitCount) incr() {
	ac.count.Add(1)
}

func (ac *awaitCount) decr() {
	if num := ac.count.Add(-1); num <= 0 {
		ac.cancel()
	}
}
