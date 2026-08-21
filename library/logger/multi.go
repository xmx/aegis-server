package logger

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
)

func NewMultiHandler(hs ...slog.Handler) *MultiHandler {
	m := new(MultiHandler)
	m.Replace(hs...)

	return m
}

type MultiHandler struct {
	holder atomic.Pointer[slog.MultiHandler]
	mutex  sync.Mutex
	order  []slog.Handler            // 保证输出顺序
	unique map[slog.Handler]struct{} // 保证不重复
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return m.load().Enabled(ctx, level)
}

func (m *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	return m.load().Handle(ctx, record)
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return m.load().WithAttrs(attrs)
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	return m.load().WithGroup(name)
}

func (m *MultiHandler) Append(hs ...slog.Handler) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.unique == nil {
		m.unique = make(map[slog.Handler]struct{}, len(hs))
	}

	for _, h := range hs {
		if h == nil {
			continue
		}
		if _, ok := m.unique[h]; ok {
			continue
		}

		m.unique[h] = struct{}{}
		m.order = append(m.order, h)
	}

	m.holder.Store(slog.NewMultiHandler(m.order...))
}

func (m *MultiHandler) Remove(hs ...slog.Handler) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, h := range hs {
		if _, ok := m.unique[h]; !ok {
			continue
		}

		delete(m.unique, h)
		for i, v := range m.order {
			if v == h {
				m.order = append(m.order[:i], m.order[i+1:]...)
				break
			}
		}
	}

	m.holder.Store(slog.NewMultiHandler(m.order...))
}

func (m *MultiHandler) Replace(hs ...slog.Handler) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	unique := make(map[slog.Handler]struct{}, len(hs))
	order := make([]slog.Handler, 0, len(hs))
	for _, h := range hs {
		if h == nil {
			continue
		}
		if _, ok := unique[h]; ok {
			continue
		}

		order = append(order, h)
		unique[h] = struct{}{}
	}

	m.order = order
	m.unique = unique
	m.holder.Store(slog.NewMultiHandler(order...))
}

func (m *MultiHandler) load() *slog.MultiHandler {
	if h := m.holder.Load(); h != nil {
		return h
	}

	return slog.NewMultiHandler()
}
