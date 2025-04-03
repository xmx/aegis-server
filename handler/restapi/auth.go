package restapi

import "github.com/xmx/ship"

func NewAuth() *Auth {
	return &Auth{}
}

type Auth struct{}

func (a *Auth) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/auth/back").GET(a.back)
	return nil
}

func (a *Auth) back(c *ship.Context) error {
	return nil
}
