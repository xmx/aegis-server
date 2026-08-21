package netutil

import (
	"net/http"
	"net/url"
	"time"
)

func NewTransport(d Dialer) *http.Transport {
	tp := &http.Transport{
		Proxy:               ProxyFromContext,
		MaxIdleConns:        5,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     2,
		IdleConnTimeout:     10 * time.Minute,
	}
	if d != nil {
		tp.DialContext = d.DialContext
	}

	return tp
}

func ProxyFromContext(r *http.Request) (*url.URL, error) {
	if cv := FromContext(r.Context()); cv != nil && cv.ProxyURL != "" {
		return url.Parse(cv.ProxyURL)
	}

	return nil, nil
}
