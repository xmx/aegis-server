package memoize

import (
	"context"
	"fmt"
	"sync"
)

func Map[K comparable, V any](call func(ctx context.Context, key K) (V, error)) Mapping[K, V] {
	return &cacheMap[K, V]{
		call:    call,
		entries: make(map[K]V, 16),
	}
}

type cacheMap[K comparable, V any] struct {
	call    func(ctx context.Context, key K) (V, error)
	mutex   sync.RWMutex
	entries map[K]V
}

func (cm *cacheMap[K, V]) Load(ctx context.Context, key K) (V, error) {
	cm.mutex.RLock()
	data, ok := cm.fastLoad(key)
	cm.mutex.RUnlock()
	if ok {
		return data, nil
	}

	return cm.slowLoad(ctx, key)
}

func (cm *cacheMap[K, V]) Forget(key K) {
	cm.mutex.Lock()
	delete(cm.entries, key)
	cm.mutex.Unlock()
}

func (cm *cacheMap[K, V]) Forgets() {
	cm.mutex.Lock()
	cm.entries = make(map[K]V, 16)
	cm.mutex.Unlock()
}

func (cm *cacheMap[K, V]) fastLoad(key K) (V, bool) {
	v, ok := cm.entries[key]
	return v, ok
}

func (cm *cacheMap[K, V]) slowLoad(ctx context.Context, key K) (V, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if data, ok := cm.fastLoad(key); ok {
		return data, nil
	}

	data, err := cm.safeCall(ctx, key)
	if err == nil {
		cm.entries[key] = data
	}

	return data, err
}

func (cm *cacheMap[K, V]) safeCall(ctx context.Context, key K) (data V, err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("panicked oncall: %v", v)
		}
	}()

	data, err = cm.call(ctx, key)

	return
}
