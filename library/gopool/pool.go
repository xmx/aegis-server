package gopool

import (
	"context"
	"sync/atomic"
)

type Pool interface {
	Go(f func()) (done <-chan struct{})
	Gos(parent context.Context, funcs ...func(parent context.Context)) (done <-chan struct{})
}

func NewPool(worker int) Pool {
	if worker <= 0 {
		worker = 1
	}
	return &pool{
		sema: make(chan struct{}, worker),
	}
}

type pool struct {
	sema chan struct{}
}

func (p *pool) Go(f func()) <-chan struct{} {
	done := make(chan struct{})
	if f == nil {
		close(done)
		return done
	}

	fn := p.warp(done, f)
	p.join(fn)

	return done
}

func (p *pool) Gos(parent context.Context, funcs ...func(parent context.Context)) <-chan struct{} {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	mon := p.newMonitor(ctx, cancel)
	fns := make([]func(), 0, len(funcs))
	for _, f := range funcs {
		if f != nil {
			fn := mon.warp(f)
			fns = append(fns, fn)
		}
	}
	if len(fns) == 0 {
		cancel()
		return ctx.Done()
	}
	for _, fn := range fns {
		p.join(fn)
	}

	return ctx.Done()
}

func (p *pool) join(f func()) {
	p.sema <- struct{}{}
	go f()
}

func (p *pool) warp(done chan struct{}, f func()) func() {
	return func() {
		defer func() {
			close(done)
			<-p.sema
		}()
		f()
	}
}

func (p *pool) newMonitor(ctx context.Context, cancel context.CancelFunc) *monitor {
	return &monitor{ctx: ctx, cancel: cancel}
}

type monitor struct {
	cnt    atomic.Int32
	ctx    context.Context
	cancel context.CancelFunc
}

func (m *monitor) warp(f func(context.Context)) func() {
	m.cnt.Add(1)
	return func() {
		defer func() {
			if num := m.cnt.Add(-1); num <= 0 {
				m.cancel()
			}
		}()
		f(m.ctx)
	}
}
