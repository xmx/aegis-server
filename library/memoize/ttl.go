package memoize

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TTL Time To Live.
func TTL[T any](call func(context.Context) (T, error), ttl time.Duration) Caching[T] {
	return &cacheTTL[T]{
		call: call,
		ttl:  ttl,
	}
}

type cacheTTL[T any] struct {
	call   func(ctx context.Context) (T, error)
	ttl    time.Duration
	mutex  sync.RWMutex
	callAt time.Time
	data   T
}

func (ct *cacheTTL[T]) Load(ctx context.Context) (T, error) {
	now := time.Now()
	ct.mutex.RLock()
	dat, valid := ct.fastLoad(now)
	ct.mutex.RUnlock()
	if valid {
		return dat, nil
	}

	return ct.slowLoad(ctx, now)
}

func (ct *cacheTTL[T]) Forget() {
	ct.mutex.Lock()
	ct.callAt = time.Time{}
	ct.mutex.Unlock()
}

func (ct *cacheTTL[T]) fastLoad(now time.Time) (T, bool) {
	data := ct.data
	valid := !ct.callAt.Add(ct.ttl).Before(now)

	return data, valid
}

func (ct *cacheTTL[T]) slowLoad(ctx context.Context, now time.Time) (T, error) {
	ct.mutex.Lock()
	defer ct.mutex.Unlock()

	if dat, valid := ct.fastLoad(now); valid {
		return dat, nil
	}

	callAt := time.Now()
	data, err := ct.safeCall(ctx)
	if err == nil {
		ct.data = data
		ct.callAt = callAt
	}

	return data, err
}

func (ct *cacheTTL[T]) safeCall(ctx context.Context) (data T, err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("panicked oncall: %v", v)
		}
	}()

	data, err = ct.call(ctx)

	return
}
