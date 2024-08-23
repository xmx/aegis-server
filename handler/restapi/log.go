package restapi

import (
	"net/http"
	"sync/atomic"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/argument/response"
	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/handler/shipx"
	"github.com/xmx/aegis-server/protocol/eventsource"
)

func NewLog(svc service.Log) shipx.Controller {
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

func (api *logAPI) Register(rt shipx.Router) error {
	auth := rt.Auth()
	auth.Route("/log/level").
		GET(api.Level).
		POST(api.SetLevel)
	auth.Route("/log/tail").GET(api.Tail)

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
	defer api.count.Add(-1)
	if cnt > api.limit {
		return ship.ErrTooManyRequests
	}

	w, r := c.ResponseWriter(), c.Request()
	sse, ok := eventsource.Accept(w, r)
	if !ok {
		return ship.ErrBadRequest
	}

	c.Warnf("进入日志查看器")
	api.svc.Attach(sse)
	defer api.svc.Leave(sse)
	<-sse.Done()
	c.Warnf("离开日志查看器")

	return nil
}
