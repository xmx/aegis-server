package memoize

import (
	"context"
	"sync"
)

type Cache[V any] interface {
	Load(context.Context) V
	Forget() V
}

func NewCache[V any](load func(context.Context) V) Cache[V] {
	return &memCache[V]{
		load: load,
	}
}

type memCache[V any] struct {
	load  func(context.Context) V
	mutex sync.RWMutex
	done  bool
	v     V
}

func (c *memCache[V]) Load(ctx context.Context) V {
	c.mutex.RLock()
	done, v := c.done, c.v
	c.mutex.RUnlock()
	if done {
		return v
	}

	return c.slowLoad(ctx)
}

func (c *memCache[V]) Forget() (v V) {
	c.mutex.Lock()
	if c.done {
		v = c.v
	}
	c.done = false
	c.mutex.Unlock()

	return
}

func (c *memCache[V]) slowLoad(ctx context.Context) V {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.done {
		return c.v
	}

	v := c.load(ctx)
	c.v, c.done = v, true

	return v
}
