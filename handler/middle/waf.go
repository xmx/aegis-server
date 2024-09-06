package middle

import (
	"log/slog"
	"time"

	"github.com/xgfone/ship/v5"
)

func WAF(log *slog.Logger) ship.Middleware {
	waf := &wafMiddle{log: log}
	return waf.middle
}

type wafMiddle struct {
	log *slog.Logger
}

func (wm *wafMiddle) middle(h ship.Handler) ship.Handler {
	return func(c *ship.Context) error {
		accessAt := time.Now()

		method, url := c.Method(), c.Request().URL
		remoteAddr, clientIP := c.RemoteAddr(), c.ClientIP()
		headers := c.ReqHeader()
		attrs := []any{
			slog.String("method", method),
			slog.String("path", url.String()),
			slog.String("remote_addr", remoteAddr),
			slog.String("client_ip", clientIP),
			slog.Any("headers", headers),
			slog.Time("access_at", accessAt),
		}

		err := h(c)
		leaveAt := time.Now()
		elapsed := leaveAt.Sub(accessAt)
		attrs = append(attrs, slog.Time("leave_at", leaveAt), slog.String("elapsed", elapsed.String()))

		if err != nil {
			attrs = append(attrs, slog.String("error", err.Error()))
			wm.log.Warn("接口访问日志", attrs...)
		} else {
			wm.log.Info("接口访问日志", attrs...)
		}

		return err
	}
}
