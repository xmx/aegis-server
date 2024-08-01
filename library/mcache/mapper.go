package mcache

import (
	"context"
	"sync"
)

type Mapper[K comparable, V, E any] interface {
	Load(context.Context, K) (V, E)
	Forget(K) (V, E)
}

func NewMapper[K comparable, V, E any](load func(context.Context, K) (V, E)) Mapper[K, V, E] {
	return &mapper[K, V, E]{
		load:  load,
		elems: make(map[K]*elemEntry[V, E], 16),
	}
}

type elemEntry[V, E any] struct {
	v V
	e E
}

type mapper[K comparable, V, E any] struct {
	load  func(context.Context, K) (V, E)
	mutex sync.RWMutex
	elems map[K]*elemEntry[V, E]
}

func (m *mapper[K, V, E]) Load(ctx context.Context, k K) (V, E) {
	m.mutex.RLock()
	ele, ok := m.elems[k]
	m.mutex.RUnlock()
	if ok {
		return ele.v, ele.e
	}

	return m.slowLoad(ctx, k)
}

func (m *mapper[K, V, E]) Forget(k K) (v V, e E) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	ele, ok := m.elems[k]
	if ok {
		delete(m.elems, k)
		v, e = ele.v, ele.e
	}

	return
}

func (m *mapper[K, V, E]) slowLoad(ctx context.Context, k K) (v V, e E) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if ele, ok := m.elems[k]; ok {
		v, e = ele.v, ele.e
		return
	}

	v, e = m.load(ctx, k)
	m.elems[k] = &elemEntry[V, E]{v: v, e: e}

	return
}
