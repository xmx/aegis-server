package restapi

import (
	"log/slog"
	"sync/atomic"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/request"
	"github.com/xmx/aegis-server/handler/shipx"
	"github.com/xmx/aegis-server/infra/profile"
)

func Log(writers profile.LogWriteCloser, level *slog.LevelVar) shipx.Register {
	return &logAPI{
		writers: writers,
		level:   level,
		limit:   5,
	}
}

type logAPI struct {
	writers profile.LogWriteCloser // 日志输出。
	level   *slog.LevelVar         // 日志级别。
	limit   int32                  // tail 日志最大个数。
	count   atomic.Int32           // tail 连接计数器。
}

func (api *logAPI) Register(rt shipx.Router) error {
	auth := rt.Auth()
	auth.Route("/log/level").POST(api.Level)
	auth.Route("/log/tail").GET(api.Tail)

	return nil
}

func (api *logAPI) Level(c *ship.Context) error {
	req := new(request.LogLevel)
	if err := c.Bind(req); err != nil {
		return err
	}

	lvl := req.Level
	level := api.level.Level()
	if err := api.level.UnmarshalText([]byte(lvl)); err == nil {
		c.Errorf("修改日志级别从 %s 修改到 %s", level, lvl)
	}

	return nil
}

func (api *logAPI) Tail(c *ship.Context) error {
	cnt := api.count.Add(1)
	defer api.count.Add(-1)
	if cnt > api.limit {
		return ship.ErrTooManyRequests
	}

	sse, err := shipx.SSE(c)
	if err != nil {
		return err
	}

	api.writers.Append(sse)
	defer api.writers.Remove(sse)
	<-sse.Done()

	return nil
}
