package memoize

import (
	"context"
	"fmt"
	"sync"
)

type Caching[T any] interface {
	Load(ctx context.Context) (T, error)
	Forget() T
}

type Mapping[K comparable, V any] interface {
	Load(ctx context.Context, key K) (V, error)
	Forget(key K)
	Forgets()
}

func Cache[T any](call func(context.Context) (T, error)) Caching[T] {
	return &cacheData[T]{
		call: call,
	}
}

type cacheData[T any] struct {
	call  func(ctx context.Context) (T, error)
	mutex sync.RWMutex
	done  bool
	data  T
}

func (cd *cacheData[T]) Load(ctx context.Context) (T, error) {
	cd.mutex.RLock()
	data, done := cd.fastLoad()
	cd.mutex.RUnlock()
	if done {
		return data, nil
	}

	return cd.slowLoad(ctx)
}

func (cd *cacheData[T]) Forget() T {
	cd.mutex.Lock()
	data := cd.data
	cd.done = false
	cd.mutex.Unlock()

	return data
}

func (cd *cacheData[T]) fastLoad() (T, bool) {
	return cd.data, cd.done
}

func (cd *cacheData[T]) slowLoad(ctx context.Context) (T, error) {
	cd.mutex.Lock()
	defer cd.mutex.Unlock()

	if data, done := cd.fastLoad(); done {
		return data, nil
	}

	data, err := cd.safeCall(ctx)
	if err == nil {
		cd.data = data
		cd.done = true
	}

	return data, err
}

func (cd *cacheData[T]) safeCall(ctx context.Context) (data T, err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("panicked oncall: %v", v)
		}
	}()

	data, err = cd.call(ctx)

	return
}
