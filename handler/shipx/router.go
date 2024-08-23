package shipx

import "github.com/xgfone/ship/v5"

type Router interface {
	Anon() *ship.RouteGroupBuilder
	Auth() *ship.RouteGroupBuilder
}

type Controller interface {
	Register(Router) error
}

func NewRouter(anon, auth *ship.RouteGroupBuilder) Router {
	return &groupRoute{
		anon: anon,
		auth: auth,
	}
}

type groupRoute struct {
	anon *ship.RouteGroupBuilder
	auth *ship.RouteGroupBuilder
}

func (g *groupRoute) Anon() *ship.RouteGroupBuilder {
	return g.anon
}

func (g *groupRoute) Auth() *ship.RouteGroupBuilder {
	return g.auth
}

type Enhancer interface {
	Describe(c *ship.Context) (name string, omit bool)
}
