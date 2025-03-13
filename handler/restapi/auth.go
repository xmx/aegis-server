package restapi

import (
	"fmt"

	"github.com/xgfone/ship/v5"
)

func NewAuth() *Auth {
	return &Auth{}
}

type Auth struct{}

func (a *Auth) Route(r *ship.RouteGroupBuilder) error {
	r.Route("/auth/back").GET(a.back)
	return nil
}

func (a *Auth) back(c *ship.Context) error {
	fmt.Println("UA：" + c.UserAgent())
	return nil
}
