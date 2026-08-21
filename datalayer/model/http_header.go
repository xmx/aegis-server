package model

import (
	"net/http"
	"strings"
)

type HTTPHeader map[string]string

func (h HTTPHeader) Canonical() HTTPHeader {
	hm := make(HTTPHeader, len(h))
	for k, v := range h {
		key := http.CanonicalHeaderKey(strings.TrimSpace(k))
		hm[key] = v
	}

	return hm
}
