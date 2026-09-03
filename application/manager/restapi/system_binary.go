package restapi

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/xmx/aegis-common/bininfo"
	"github.com/xmx/aegis-server/application/manager/service"
)

type SystemBinary struct {
	svc *service.Agent
}

func NewSystemBinary() *SystemBinary {
	return &SystemBinary{}
}

func (sys *SystemBinary) RegisterRoute(r *echo.Group) error {
	r.GET("/system/buildinfo", sys.buildinfo)

	return nil
}

func (sys *SystemBinary) buildinfo(c *echo.Context) error {
	ret := bininfo.Get()
	return c.JSON(http.StatusOK, ret)
}
