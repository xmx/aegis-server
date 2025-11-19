package restapi

import (
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-control/datalayer/model"
	"github.com/xmx/aegis-server/application/expose/service"
)

func NewSetting(svc *service.Setting) *Setting {
	return &Setting{
		svc: svc,
	}
}

type Setting struct {
	svc *service.Setting
}

func (set *Setting) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/setting").GET(set.get).POST(set.upsert)
	return nil
}

func (set *Setting) get(c *ship.Context) error {
	ctx := c.Request().Context()
	cfg, _ := set.svc.Get(ctx)

	return c.JSON(http.StatusOK, cfg)
}

func (set *Setting) upsert(c *ship.Context) error {
	req := new(model.SettingData)
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	return set.svc.Upsert(ctx, req)
}
