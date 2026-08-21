package shipx

import "github.com/xgfone/ship/v5"

type RouteRegister interface {
	RegisterRoute(*ship.RouteGroupBuilder) error
}

func RegisterRoutes(r *ship.RouteGroupBuilder, rbs []RouteRegister) error {
	for _, rb := range rbs {
		if err := rb.RegisterRoute(r); err != nil {
			return err
		}
	}

	return nil
}
