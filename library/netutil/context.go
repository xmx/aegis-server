package netutil

import "context"

type ContextValue struct {
	DialTags []string
	ProxyURL string
}

type contextKey struct {
	name string
}

var defaultContextKey = &contextKey{name: "netutil-context-key"}

func FromContext(ctx context.Context) *ContextValue {
	val, _ := ctx.Value(defaultContextKey).(*ContextValue)
	return val
}

func WithProxyURL(parent context.Context, proxyURL string) context.Context {
	if val, ok := parent.Value(defaultContextKey).(*ContextValue); ok {
		val.ProxyURL = proxyURL
		return parent
	}

	val := &ContextValue{ProxyURL: proxyURL}

	return context.WithValue(parent, defaultContextKey, val)
}

func WithDialTags(parent context.Context, tags []string) context.Context {
	if val, ok := parent.Value(defaultContextKey).(*ContextValue); ok {
		val.DialTags = tags
		return parent
	}

	val := &ContextValue{DialTags: tags}

	return context.WithValue(parent, defaultContextKey, val)
}
