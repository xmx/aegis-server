package shipx

import "github.com/xmx/ship"

type RouteRegister interface {
	RegisterRoute(r *ship.RouteGroupBuilder) error
}

func RegisterRoutes(r *ship.RouteGroupBuilder, rts []RouteRegister) error {
	for _, rt := range rts {
		if err := rt.RegisterRoute(r); err != nil {
			return err
		}
	}
	return nil
}
