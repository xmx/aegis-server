package memoize

import (
	"context"
	"sync"
)

type Cache2[K, V any] interface {
	Load(context.Context) (K, V)
	Forget() (K, V)
}

func NewCache2[K, V any](load func(context.Context) (K, V)) Cache2[K, V] {
	return &cache2[K, V]{
		load: load,
	}
}

type cache2[K, V any] struct {
	load  func(context.Context) (K, V)
	mutex sync.RWMutex
	done  bool
	k     K
	v     V
}

func (c *cache2[K, V]) Load(ctx context.Context) (K, V) {
	c.mutex.RLock()
	done := c.done
	k, v := c.k, c.v
	c.mutex.RUnlock()
	if done {
		return k, v
	}

	return c.slowLoad(ctx)
}

func (c *cache2[K, V]) Forget() (k K, v V) {
	c.mutex.Lock()
	if c.done {
		k = c.k
		v = c.v
	}
	c.done = false
	c.mutex.Unlock()

	return
}

func (c *cache2[K, V]) slowLoad(ctx context.Context) (K, V) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.done {
		return c.k, c.v
	}

	k, v := c.load(ctx)
	c.k, c.v = k, v
	c.done = true

	return k, v
}
