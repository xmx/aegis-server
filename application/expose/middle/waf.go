package middle

import (
	"context"
	"io"
	"log/slog"
	"net"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-control/datalayer/model"
)

func NewWAF(writeLog func(context.Context, *model.Oplog) error) ship.Middleware {
	waf := &wafMiddle{writeLog: writeLog}
	return waf.middle
}

type wafMiddle struct {
	writeLog func(context.Context, *model.Oplog) error
}

func (wm *wafMiddle) middle(h ship.Handler) ship.Handler {
	return func(c *ship.Context) error {
		remoteAddr, clientIP := c.RemoteAddr(), c.ClientIP()
		directIP, _, _ := net.SplitHostPort(remoteAddr)
		if directIP == "" {
			directIP = remoteAddr
		}
		host, method := c.Host(), c.Method()
		req := c.Request()
		reqURL := req.URL

		attrs := []any{
			slog.String("client_ip", clientIP),
			slog.String("remote_addr", remoteAddr),
			slog.String("method", method),
			slog.String("host", host),
			slog.String("path", reqURL.Path),
		}

		err := h(c)
		if err != nil {
			attrs = append(attrs, "error", err)
			c.Warnf("访问接口出错", attrs...)
		} else {
			c.Infof("访问接口", attrs...)
		}

		return err
	}
}

func (wm *wafMiddle) newRecordBody(body io.ReadCloser, size int) *recordBody {
	return &recordBody{
		body: body,
		data: make([]byte, size),
	}
}

type recordBody struct {
	body io.ReadCloser
	data []byte
	pos  int
}

func (rb *recordBody) Read(p []byte) (int, error) {
	n, err := rb.body.Read(p)
	i := copy(rb.data[rb.pos:], p[:n])
	rb.pos += i
	return n, err
}

func (rb *recordBody) Close() error {
	return rb.body.Close()
}

func (rb *recordBody) Data() []byte {
	return rb.data[:rb.pos]
}
