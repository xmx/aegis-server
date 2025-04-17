package queue

import (
	"sync"
	"time"
)

type overstock[T any] struct {
	items    []T
	mutex    sync.Mutex
	full     int
	timeout  time.Duration
	callback func([]T)
	timer    *time.Timer
}

func NewOverstock[T any](full int, timeout time.Duration, callback func([]T)) Queuer[T] {
	if full <= 0 {
		full = 1
	}
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	if callback == nil {
		callback = func([]T) {}
	}
	ost := &overstock[T]{
		items:    make([]T, 0, full),
		full:     full,
		timeout:  timeout,
		callback: callback,
	}
	ost.timer = time.AfterFunc(ost.timeout, ost.timeoutTrigger)

	return ost
}

func (ost *overstock[T]) Enqueue(item T) {
	ost.mutex.Lock()
	ost.items = append(ost.items, item)
	items := ost.items
	fulled := len(ost.items) >= ost.full
	if fulled {
		ost.items = make([]T, 0, ost.full)
		ost.timer.Reset(ost.timeout)
	}
	ost.mutex.Unlock()

	if fulled {
		ost.callback(items)
	}
}

func (ost *overstock[T]) Dequeue() (T, bool) {
	ost.mutex.Lock()
	defer ost.mutex.Unlock()

	if len(ost.items) == 0 {
		var zero T
		return zero, false
	}

	item := ost.items[0]
	ost.items = ost.items[1:]
	return item, true
}

func (ost *overstock[T]) timeoutTrigger() {
	ost.mutex.Lock()
	ost.timer.Reset(ost.timeout)
	items := ost.items
	ost.items = make([]T, 0, ost.full)
	ost.mutex.Unlock()

	if len(items) > 0 {
		ost.callback(items)
	}
}
