package memoize

import (
	"context"
	"sync"
	"time"
)

type TTLCache2[V, E any] interface {
	Load(context.Context) (V, E)

	// Forget 清除内存中缓存的数据，不返回数据。
	Forget()
}

func NewTTL2[V, E any](fn func(context.Context) (V, E), ttl time.Duration) TTLCache2[V, E] {
	if ttl < 0 {
		ttl = 0
	}

	return &cacheTTL2[V, E]{
		fn: fn,
		du: ttl,
	}
}

type cacheTTL2[V, E any] struct {
	fn  func(context.Context) (V, E)
	du  time.Duration
	mu  sync.RWMutex
	exp time.Time
	ent *entry2[V, E]
}

func (tc *cacheTTL2[V, E]) Load(ctx context.Context) (V, E) {
	now := time.Now()
	tc.mu.RLock()
	ent := tc.loadEntry(now)
	tc.mu.RUnlock()
	if ent != nil {
		return ent.load()
	}

	return tc.slowLoad(ctx, now)
}

func (tc *cacheTTL2[V, E]) Forget() {
	tc.mu.Lock()
	tc.exp = time.Time{}
	tc.ent = nil
	tc.mu.Unlock()
}

func (tc *cacheTTL2[V, E]) loadEntry(now time.Time) *entry2[V, E] {
	if tc.ent != nil && !now.After(tc.exp) {
		return tc.ent
	}

	return nil
}

func (tc *cacheTTL2[V, E]) slowLoad(ctx context.Context, now time.Time) (V, E) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if ent := tc.loadEntry(now); ent != nil {
		return ent.load()
	}

	v, e := tc.fn(ctx)
	tc.ent = &entry2[V, E]{v: v, e: e}
	tc.exp = now.Add(tc.du)

	return v, e
}
