package middle

import (
	"net/http"

	"github.com/xgfone/ship/v5"
)

type Allower interface {
	Allowed(*http.Request) bool
}

func NewFirewall(allow Allower) ship.Middleware {
	fw := &firewall{allow: allow}
	return fw.handle
}

type firewall struct {
	allow Allower
}

func (f *firewall) handle(h ship.Handler) ship.Handler {
	return func(c *ship.Context) error {
		r := c.Request()
		if f.allow.Allowed(r) {
			return h(c)
		}

		return ship.ErrForbidden
	}
}
