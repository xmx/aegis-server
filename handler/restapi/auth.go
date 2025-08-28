package restapi

import (
	"fmt"

	"github.com/xgfone/ship/v5"
)

func NewAuth() *Auth {
	return &Auth{}
}

type Auth struct{}

func (a *Auth) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/auth/back").GET(a.back)
	return nil
}

func (a *Auth) back(c *ship.Context) error {
	req := new(User)
	err := c.BindQuery(req)
	fmt.Println(err)
	return nil
}

type User struct {
	Passwd string `json:"passwd" validate:"password"`
}
