package restapi

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lmittmann/tint"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-common/library/useragent"
	"github.com/xmx/aegis-common/logger"
	"github.com/xmx/aegis-server/contract/request"
	"github.com/xmx/aegis-server/protocol/eventsource"
)

func NewLog(handler logger.Handler) *Log {
	return &Log{
		handler: handler,
		upg: &websocket.Upgrader{
			HandshakeTimeout:  30 * time.Second,
			CheckOrigin:       func(*http.Request) bool { return true },
			EnableCompression: true,
		},
	}
}

type Log struct {
	handler logger.Handler
	upg     *websocket.Upgrader
}

func (l *Log) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/log/watch").GET(l.watch)

	return nil
}

// watch 观测日志。
//
//goland:noinspection GoUnhandledErrorResult
func (l *Log) watch(c *ship.Context) error {
	req := new(request.LogWatch)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	parent := c.Request().Context()
	var writer doneWriteCloser
	w, r := c.ResponseWriter(), c.Request()
	if c.IsWebSocket() {
		conn, err := l.upg.Upgrade(w, r, nil)
		if err != nil {
			return err
		}
		ctx, cancel := context.WithCancel(parent)
		writer = &socketLog{conn: conn, ctx: ctx, cancel: cancel}
	} else if sse := eventsource.Accept(w, r); sse != nil {
		writer = sse
	} else {
		w.Header().Set(ship.HeaderContentType, ship.MIMETextPlainCharsetUTF8)
		w.WriteHeader(http.StatusOK)
		ctx, cancel := context.WithCancel(parent)
		writer = &streamableHTTP{wrt: w, ctx: ctx, cancel: cancel}
	}
	defer writer.Close()

	level := req.LevelVar()
	opts := &slog.HandlerOptions{AddSource: true, Level: level}
	var handler slog.Handler

	userAgent := c.UserAgent()
	if req.Format == "" && useragent.IsSupportedANSI(userAgent) {
		handler = tint.NewHandler(writer, &tint.Options{
			AddSource:   opts.AddSource,
			Level:       opts.Level,
			ReplaceAttr: opts.ReplaceAttr,
			TimeFormat:  time.RFC3339,
		})
	} else if req.JSONFormat() {
		handler = slog.NewJSONHandler(writer, opts)
	} else {
		handler = slog.NewTextHandler(writer, opts)
	}
	l.handler.Attach(handler)
	defer l.handler.Detach(handler)
	c.Infof("HELLO")
	<-writer.Done()

	return nil
}

type streamableHTTP struct {
	wrt    http.ResponseWriter
	ctx    context.Context
	cancel context.CancelFunc
}

func (sh *streamableHTTP) Write(p []byte) (int, error) {
	n, err := sh.wrt.Write(p)
	if f, ok := sh.wrt.(http.Flusher); ok {
		f.Flush()
	}

	return n, err
}

func (sh *streamableHTTP) Close() error {
	sh.cancel()
	return nil
}

func (sh *streamableHTTP) Done() <-chan struct{} {
	return sh.ctx.Done()
}

type doneWriteCloser interface {
	io.WriteCloser
	Done() <-chan struct{}
}

type socketLog struct {
	conn   *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

func (sl *socketLog) Write(p []byte) (int, error) {
	n := len(p)
	err := sl.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (sl *socketLog) Close() error {
	sl.cancel()
	return sl.conn.Close()
}

func (sl *socketLog) Done() <-chan struct{} {
	return sl.ctx.Done()
}
