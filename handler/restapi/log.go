package restapi

import (
	"net/http"
	"sync/atomic"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/errcode"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/argument/response"
	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/handler/shipx"
	"github.com/xmx/aegis-server/protocol/eventsource"
)

func NewLog(svc service.Log) shipx.Router {
	return &logAPI{
		svc:   svc,
		limit: 5,
	}
}

type logAPI struct {
	svc   service.Log
	limit int32        // tail 日志最大个数。
	count atomic.Int32 // tail 连接计数器。
}

func (api *logAPI) Route(r *ship.RouteGroupBuilder) error {
	r.Route("/log/level").
		GET(api.Level).
		POST(api.SetLevel)
	r.Route("/sse/log/tail").GET(api.Tail)

	return nil
}

func (api *logAPI) Level(c *ship.Context) error {
	level := api.svc.Level()
	ret := &response.LogLevel{Level: level}
	return c.JSON(http.StatusOK, ret)
}

func (api *logAPI) SetLevel(c *ship.Context) error {
	req := new(request.LogLevel)
	if err := c.Bind(req); err != nil {
		return err
	}

	return api.svc.SetLevel(req.Level)
}

func (api *logAPI) Tail(c *ship.Context) error {
	cnt := api.count.Add(1)
	if lim := api.limit; cnt > lim {
		api.count.Add(-1)
		return errcode.ErrTooManyRequests
	}

	defer func() {
		num := api.count.Add(-1)
		c.Infof("离开日志查看器，当前还有 %d 个查看窗口", num)
	}()

	w, r := c.ResponseWriter(), c.Request()
	sse := eventsource.Accept(w, r)
	if sse == nil {
		c.Warnf("不是 Server-Sent Events 连接")
		return errcode.ErrServerSentEvents
	}

	api.svc.Attach(sse)
	defer api.svc.Leave(sse)

	c.Warnf("进入日志查看器，当前共有 %d 个查看窗口", cnt)
	<-sse.Done()

	return nil
}
