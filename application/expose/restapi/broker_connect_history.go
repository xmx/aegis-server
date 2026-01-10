package restapi

import (
	"net/http"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/application/expose/request"
	"github.com/xmx/aegis-server/application/expose/service"
)

func NewBrokerConnectHistory(svc *service.BrokerConnectHistory) *BrokerConnectHistory {
	return &BrokerConnectHistory{
		svc: svc,
	}
}

type BrokerConnectHistory struct {
	svc *service.BrokerConnectHistory
}

func (bch *BrokerConnectHistory) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/agent/connect/histories").GET(bch.page)

	return nil
}

func (bch *BrokerConnectHistory) page(c *ship.Context) error {
	req := new(request.PageKeywords)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	ret, err := bch.svc.Page(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}
