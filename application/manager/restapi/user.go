package restapi

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/xmx/aegis-server/application/manager/request"
	"github.com/xmx/aegis-server/application/manager/service"
)

type User struct {
	svc *service.User
}

func NewUser(svc *service.User) *User {
	return &User{
		svc: svc,
	}
}

func (usr *User) RegisterRoute(r *echo.Group) error {
	r.GET("/users", usr.page)

	return nil
}

func (usr *User) page(c *echo.Context) error {
	req := new(request.Pages)
	if err := c.Bind(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	ret, err := usr.svc.Page(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}
