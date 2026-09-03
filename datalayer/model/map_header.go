package model

import (
	"net/http"
	"strings"
)

type MapHeader map[string]string

func (h MapHeader) Canonical() MapHeader {
	hm := make(MapHeader, len(h))
	for k, v := range h {
		key := http.CanonicalHeaderKey(strings.TrimSpace(k))
		hm[key] = v
	}

	return hm
}
