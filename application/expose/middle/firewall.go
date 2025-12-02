package middle

import (
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/application/expose/firewalld"
)

func NewFirewall(fw *firewalld.Firewalld) ship.Middleware {
	f := &firewall{fw: fw}
	return f.handle
}

type firewall struct {
	fw *firewalld.Firewalld
}

func (f *firewall) handle(h ship.Handler) ship.Handler {
	return func(c *ship.Context) error {
		if fw := f.fw; fw == nil || fw.Allow(c.Request()) {
			return h(c)
		}

		return ship.ErrForbidden.Newf("禁止访问")
	}
}
