package restapi

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/xmx/aegis-common/psutil"
)

type SystemService struct{}

func NewSystemService() *SystemService {
	return &SystemService{}
}

func (ps *SystemService) RegisterRoute(r *echo.Group) error {
	r.GET("/system/services", ps.list)

	return nil
}

func (ps *SystemService) list(c *echo.Context) error {
	ret, err := psutil.Services()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}
