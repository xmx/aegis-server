package restapi

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/xmx/aegis-common/psutil"
	"github.com/xmx/aegis-server/application/manager/request"
)

type SystemProcess struct{}

func NewSystemProcess() *SystemProcess {
	return &SystemProcess{}
}

func (ps *SystemProcess) RegisterRoute(r *echo.Group) error {
	r.GET("/system/process/samples", ps.samples)
	r.GET("/system/process/details", ps.details)

	return nil
}

func (ps *SystemProcess) samples(c *echo.Context) error {
	ctx := c.Request().Context()
	ret, err := psutil.ProcessSamples(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (ps *SystemProcess) details(c *echo.Context) error {
	req := new(request.PID)
	if err := c.Bind(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	ret, err := psutil.Process(ctx, req.PID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}
