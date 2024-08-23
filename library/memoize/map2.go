package memoize

import (
	"context"
	"sync"
)

type Map2[K comparable, V, E any] interface {
	Load(context.Context, K) (V, E)
	Forget(K) (V, E)
	Forgets()
}

func NewMap2[K comparable, V, E any](fn func(context.Context, K) (V, E)) Map2[K, V, E] {
	return &cacheMap2[K, V, E]{
		fn: fn,
		hm: make(map[K]*entry2[V, E], 8),
	}
}

type cacheMap2[K comparable, V, E any] struct {
	fn func(context.Context, K) (V, E)
	mu sync.RWMutex
	hm map[K]*entry2[V, E]
}

func (cm *cacheMap2[K, V, E]) Forget(k K) (v V, e E) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if ent, ok := cm.hm[k]; ok {
		v, e = ent.load()
		delete(cm.hm, k)
	}

	return
}

func (cm *cacheMap2[K, V, E]) Forgets() {
	cm.mu.Lock()
	cm.hm = make(map[K]*entry2[V, E], 8)
	cm.mu.Unlock()
}

func (cm *cacheMap2[K, V, E]) Load(ctx context.Context, k K) (V, E) {
	cm.mu.RLock()
	ent, ok := cm.hm[k]
	cm.mu.RUnlock()
	if ok {
		return ent.load()
	}

	return cm.slowLoad(ctx, k)
}

func (cm *cacheMap2[K, V, E]) slowLoad(ctx context.Context, k K) (V, E) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if ent, ok := cm.hm[k]; ok {
		return ent.load()
	}

	v, e := cm.fn(ctx, k)
	cm.hm[k] = &entry2[V, E]{v: v, e: e}

	return v, e
}
