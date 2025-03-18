package restapi

import "github.com/xmx/ship"

func NewAuth() *Auth {
	return &Auth{}
}

type Auth struct{}

func (a *Auth) Route(r *ship.RouteGroupBuilder) error {
	r.Route("/auth/back").GET(a.back)
	return nil
}

func (a *Auth) back(c *ship.Context) error {
	type request struct {
		Name string `json:"name" query:"name" validate:"required,lte=10"`
	}
	req := new(request)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	return nil
}
