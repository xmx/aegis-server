package memoize

import (
	"context"
	"sync"
)

type Cache2[V, E any] interface {
	Load(context.Context) (V, E)
	Forget() (V, E)
}

func NewCache2[V, E any](fn func(context.Context) (V, E)) Cache2[V, E] {
	return &cache2[V, E]{
		fn: fn,
	}
}

type cache2[V, E any] struct {
	fn  func(context.Context) (V, E)
	mu  sync.RWMutex
	ent *entry2[V, E]
}

func (ch *cache2[V, E]) Load(ctx context.Context) (V, E) {
	ch.mu.RLock()
	ent := ch.ent
	ch.mu.RUnlock()
	if ent != nil {
		return ent.load()
	}

	return ch.slowLoad(ctx)
}

func (ch *cache2[V, E]) Forget() (v V, e E) {
	ch.mu.Lock()
	if ent := ch.ent; ent != nil {
		v, e = ent.load()
		ch.ent = nil
	}
	ch.mu.Unlock()

	return
}

func (ch *cache2[V, E]) slowLoad(ctx context.Context) (V, E) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	if ent := ch.ent; ent != nil {
		return ent.load()
	}

	v, e := ch.fn(ctx)
	ch.ent = &entry2[V, E]{v: v, e: e}

	return v, e
}
