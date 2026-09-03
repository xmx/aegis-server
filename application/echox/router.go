package echox

import (
	"github.com/labstack/echo/v5"
)

type RouteRegister interface {
	RegisterRoute(r *echo.Group) error
}

func RegisterRoutes(r *echo.Group, rbs []RouteRegister) error {
	for _, rb := range rbs {
		if err := rb.RegisterRoute(r); err != nil {
			return err
		}
	}

	return nil
}
