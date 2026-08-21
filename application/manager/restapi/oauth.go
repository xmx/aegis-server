package restapi

import (
	"net/http"
	"net/url"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/application/errcode"
	"github.com/xmx/aegis-server/application/manager/request"
	"github.com/xmx/aegis-server/application/manager/response"
	"github.com/xmx/aegis-server/application/manager/service"
)

type OAuth struct {
	svc *service.OAuth
}

func NewOAuth(svc *service.OAuth) *OAuth {
	return &OAuth{
		svc: svc,
	}
}

func (ath *OAuth) github(c *ship.Context) error {
	req := new(request.OAuthCode)
	if err := c.Bind(req); err != nil {
		return err
	}

	return nil
}

func (ath *OAuth) provider(c *ship.Context) error {
	req := new(request.OAuthProvider)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	cfg, err := ath.svc.Provider(ctx, req.Provider)
	if err != nil {
		return err
	}
	if !cfg.Enabled {
		return errcode.ErrForbiddenOAuthProvider
	}

	ret := &response.OAuthProvider{
		Provider: cfg.Provider,
		ClientID: cfg.ClientID,
		Scopes:   cfg.Scopes,
	}
	if num := len(cfg.RedirectURIs); num == 1 {
		ret.RedirectURI = cfg.RedirectURIs[0]
	} else if num > 1 {
		origin := req.Origin
		if origin == "" {
			origin = c.GetReqHeader(ship.HeaderOrigin)
		}
		if hope, _ := url.Parse(origin); hope != nil {
			for _, uri := range cfg.RedirectURIs {
				pu, _ := url.Parse(uri)
				if pu == nil {
					continue
				}
				if pu.Scheme == hope.Scheme && pu.Host == hope.Host {
					ret.RedirectURI = uri
					break
				}
			}

			if ret.RedirectURI == "" {
				ret.RedirectURI = cfg.RedirectURIs[0]
			}
		}
	}

	return c.Respond(http.StatusOK, ret)
}
