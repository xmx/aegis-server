package restapi

import (
	"net/http"
	"net/url"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/application/manager/request"
	"github.com/xmx/aegis-server/application/manager/response"
	"github.com/xmx/aegis-server/application/manager/service"
	"github.com/xmx/aegis-server/datalayer/model"
	"golang.org/x/oauth2/github"
)

type OAuth struct {
	svc *service.OAuth
}

func NewOAuth(svc *service.OAuth) *OAuth {
	return &OAuth{
		svc: svc,
	}
}

func (ath *OAuth) RegisterRoute(r *ship.RouteGroupBuilder) error {
	r.Route("/oauth/provider").GET(ath.provider)
	r.Route("/oauth/github").POST(ath.github)

	return nil
}

func (ath *OAuth) github(c *ship.Context) error {
	req := new(request.OAuthCode)
	if err := c.Bind(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	ret, err := ath.svc.GitHub(ctx, req.Code, req.RedirectURI)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ret)
}

func (ath *OAuth) provider(c *ship.Context) error {
	req := new(request.OAuthProvider)
	if err := c.BindQuery(req); err != nil {
		return err
	}

	provider, origin := req.Provider, req.Origin
	if origin == "" {
		origin = c.GetReqHeader(ship.HeaderOrigin)
	}

	ctx := c.Request().Context()
	cfg, err := ath.svc.Provider(ctx, req.Provider)
	if err != nil {
		return err
	}

	redirectURI := ath.matchRedirectURI(cfg.RedirectURIs, origin)
	ret := &response.OAuthProvider{
		Provider:    cfg.Provider,
		ClientID:    cfg.ClientID,
		Scopes:      cfg.Scopes,
		RedirectURI: redirectURI,
	}
	if provider == model.OAuthGithub { // 暂时仅支持 GitHub
		ret.AuthURL = github.Endpoint.AuthURL
	}

	return c.JSON(http.StatusOK, ret)
}

func (*OAuth) matchRedirectURI(redirectURIs []string, origin string) string {
	if len(redirectURIs) == 0 {
		return ""
	}
	pu, _ := url.Parse(origin)
	if pu == nil {
		return ""
	}

	for _, uri := range redirectURIs {
		pr, err := url.Parse(uri)
		if err != nil {
			continue
		}

		if pr.Scheme == pu.Scheme && pr.Host == pu.Host {
			return uri
		}
	}

	return redirectURIs[0]
}
