package gopool

import (
	"context"
	"sync/atomic"
)

type Pool interface {
	Go(func()) context.Context
	Gos(context.Context, ...func(context.Context)) context.Context
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

func (p *pool) Go(f func()) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	if f == nil {
		cancel()
		return ctx
	}

	fn := p.warpFunc(cancel, f)
	p.join(fn)

	return ctx
}

func (p *pool) Gos(parent context.Context, funcs ...func(ctx context.Context)) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	mon := p.newMonitor(ctx, cancel)
	fns := make([]func(), 0, len(funcs))
	for _, f := range funcs {
		if f == nil {
			continue
		}
		fn := mon.warpFunc(f)
		fns = append(fns, fn)
	}
	if len(fns) == 0 {
		cancel()
		return ctx
	}

	for _, fn := range fns {
		p.join(fn)
	}

	return ctx
}

func (p *pool) join(fn func()) {
	select {
	case p.sema <- struct{}{}:
		go p.call(fn)
	default:
		select {
		case p.queue <- fn:
		case p.sema <- struct{}{}:
			go p.call(fn)
		}
	}
}

func (p *pool) call(fn func()) {
	defer func() { <-p.sema }()

	fn()
	for over := false; !over; {
		select {
		case fn = <-p.queue:
			fn()
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
