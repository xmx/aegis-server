package middle

import (
	"context"
	"io"
	"log/slog"
	"net"
	"time"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/datalayer/model"
)

func WAF(writeLog func(context.Context, *model.Oplog) error) ship.Middleware {
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
		body := wm.newRecordBody(req.Body, 4096)
		req.Body = body

		var err error
		accessedAt := time.Now()
		oplog := &model.Oplog{
			// Name:       "",
			Host:       host,
			Method:     method,
			Path:       reqURL.Path,
			Query:      reqURL.Query(),
			Header:     req.Header,
			ClientIP:   clientIP,
			DirectIP:   directIP,
			AccessedAt: accessedAt,
		}
		defer func() {
			var failed bool
			if val := recover(); val != nil {
				failed = true
			}
			if err != nil {
				failed = true
				oplog.Reason = err.Error()
			}
			oplog.FinishedAt = time.Now()
			oplog.Succeed = !failed
			oplog.Body = body.Data()
			attr := slog.Any("oplog", oplog)

			if fn := wm.writeLog; fn != nil {
				background := context.Background()
				ctx, cancel := context.WithTimeout(background, 5*time.Second)
				exx := wm.writeLog(ctx, oplog)
				cancel()
				if exx != nil {
					c.Errorf("保存访问日志出错", attr, slog.Any("error", exx))
				}
			}
			if failed {
				c.Warnf("接口访问", attr)
			} else {
				c.Infof("接口访问", attr)
			}
		}()

		err = h(c)

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
