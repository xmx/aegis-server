package netutil

import (
	"context"
	"net"
	"time"
)

type Dialer interface {
	DialContext(ctx context.Context, network, addr string) (net.Conn, error)
}

type DialMatcher interface {
	MatchDialer(ctx context.Context, tags []string) (Dialer, error)
}

type tagDialer struct {
	matcher DialMatcher
	sysdial Dialer
}

func NewTagDialer(matcher DialMatcher) Dialer {
	return &tagDialer{
		matcher: matcher,
		sysdial: &net.Dialer{Timeout: 30 * time.Second},
	}
}

func (td *tagDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	if conn, err := td.byTags(ctx, network, addr); conn != nil || err != nil {
		return conn, err
	}

	return td.sysdial.DialContext(ctx, network, addr)
}

func (td *tagDialer) byTags(ctx context.Context, network, addr string) (net.Conn, error) {
	if td.matcher == nil {
		return nil, nil
	}

	cv := FromContext(ctx)
	if cv == nil {
		return nil, nil
	}

	d, err := td.matcher.MatchDialer(ctx, cv.DialTags)
	if err != nil {
		return nil, err
	} else if d == nil {
		return nil, nil
	}

	return d.DialContext(ctx, network, addr)
}
