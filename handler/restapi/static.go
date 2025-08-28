package restapi

import "github.com/xgfone/ship/v5"

func NewStatic(sites map[string]string) *Static {
	return &Static{sites: sites}
}

type Static struct {
	sites map[string]string // PATH:DIR
}

func (sta *Static) RegisterRoute(r *ship.RouteGroupBuilder) error {
	for path, dir := range sta.sites {
		if path == "" || dir == "" {
			continue
		}
		r.Route(path).Static(dir)
	}

	return nil
}
