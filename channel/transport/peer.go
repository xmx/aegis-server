package transport

import "context"

type Peer interface {
	// ID 节点全局唯一 ID。
	ID() string

	// Muxer 底层多路通道。
	Muxer() Muxer
}

func WithValue(ctx context.Context, p Peer) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if p == nil {
		return ctx
	}

	return context.WithValue(ctx, defaultContextKey, p)
}

func FromContext(ctx context.Context) Peer {
	if ctx == nil {
		return nil
	}

	if p, ok := ctx.Value(defaultContextKey).(Peer); ok {
		return p
	}

	return nil
}

var defaultContextKey = contextKey{}

type contextKey struct{}

func (contextKey) String() string {
	return "transport-peer-context-key"
}
