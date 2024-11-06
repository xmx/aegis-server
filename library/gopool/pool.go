package gopool

import (
	"context"
	"sync/atomic"
)

type Pool interface {
	Go(fn func()) context.Context
	Gos(ctx context.Context, fns ...func(ctx context.Context)) context.Context
}

func NewPool(worker, queue int) Pool {
	if worker < 1 {
		worker = 1
	}
	if queue < 0 {
		queue = 0
	}

	return &pool{
		sema:  make(chan struct{}, worker),
		queue: make(chan func(), queue),
	}
}

type pool struct {
	sema  chan struct{}
	queue chan func()
}

func (p *pool) Go(fn func()) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	if fn == nil {
		cancel()
		return ctx
	}

	fun := p.warpFunc(cancel, fn)
	p.join(fun)

	return ctx
}

func (p *pool) Gos(parent context.Context, fns ...func(ctx context.Context)) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	mon := p.newMonitor(ctx, cancel)
	funs := make([]func(), 0, len(fns))
	for _, fn := range fns {
		if fn == nil {
			continue
		}
		fun := mon.warpFunc(fn)
		funs = append(funs, fun)
	}
	if len(funs) == 0 {
		cancel()
		return ctx
	}

	for _, fun := range funs {
		p.join(fun)
	}

	return ctx
}

func (p *pool) join(fun func()) {
	select {
	case p.sema <- struct{}{}:
		go p.call(fun)
	default:
		select {
		case p.queue <- fun:
		case p.sema <- struct{}{}:
			go p.call(fun)
		}
	}
}

func (p *pool) call(fun func()) {
	defer func() { <-p.sema }()

	fun()
	for over := false; !over; {
		select {
		case fun = <-p.queue:
			fun()
		default:
			over = true
		}
	}
}

func (p *pool) warpFunc(cancel context.CancelFunc, fn func()) func() {
	return func() {
		defer func() {
			recover()
			cancel()
		}()
		fn()
	}
}

func (p *pool) newMonitor(ctx context.Context, cancel context.CancelFunc) *monitor {
	return &monitor{
		ctx:    ctx,
		cancel: cancel,
	}
}

type monitor struct {
	count  atomic.Int32
	ctx    context.Context
	cancel context.CancelFunc
}

func (m *monitor) warpFunc(fn func(ctx context.Context)) func() {
	m.count.Add(1)
	return func() {
		defer func() {
			recover()
			if num := m.count.Add(-1); num <= 0 {
				m.cancel()
			}
		}()
		fn(m.ctx)
	}
}
