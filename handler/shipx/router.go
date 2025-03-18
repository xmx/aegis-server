package shipx

import "github.com/xmx/ship"

type Router interface {
	Route(r *ship.RouteGroupBuilder) error
}

func BindRouters(r *ship.RouteGroupBuilder, routes []Router) error {
	for _, route := range routes {
		if route == nil {
			continue
		}
		if err := route.Route(r); err != nil {
			return err
		}
	}

	return nil
}
