package restapi

import (
	"net/http"
	"sync/atomic"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/errcode"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/argument/response"
	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/protocol/eventsource"
)

func NewLog(svc *service.Log) *Log {
	return &Log{
		svc:   svc,
		limit: 5,
	}
}

type Log struct {
	svc   *service.Log
	limit int32        // tail 日志最大个数。
	count atomic.Int32 // tail 连接计数器。
}

func (l *Log) Route(r *ship.RouteGroupBuilder) error {
	r.Route("/log/level").
		GET(l.level).
		POST(l.setLevel)
	r.Route("/sse/log/tail").GET(l.tail)

	return nil
}

func (l *Log) level(c *ship.Context) error {
	lvl := l.svc.Level()
	ret := &response.LogLevel{Level: lvl}
	return c.JSON(http.StatusOK, ret)
}

func (l *Log) setLevel(c *ship.Context) error {
	req := new(request.LogLevel)
	if err := c.Bind(req); err != nil {
		return err
	}

	return l.svc.SetLevel(req.Level)
}

func (l *Log) tail(c *ship.Context) error {
	cnt := l.count.Add(1)
	if lim := l.limit; cnt > lim {
		l.count.Add(-1)
		return errcode.ErrTooManyRequests
	}

	defer func() {
		num := l.count.Add(-1)
		c.Infof("离开日志查看器，当前还有 %d 个查看窗口", num)
	}()

	w, r := c.ResponseWriter(), c.Request()
	sse := eventsource.Accept(w, r)
	if sse == nil {
		c.Warnf("不是 Server-Sent Events 连接")
		return errcode.ErrServerSentEvents
	}

	l.svc.Attach(sse)
	defer l.svc.Detach(sse)

	c.Warnf("进入日志查看器，当前共有 %d 个查看窗口", cnt)
	<-sse.Done()

	return nil
}
