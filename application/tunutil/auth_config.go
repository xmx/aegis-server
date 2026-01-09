package tunutil

import (
	"context"

	"github.com/xmx/aegis-server/application/serverd"
	"github.com/xmx/aegis-server/config"
)

func NewAuthConfig(cfg config.Database) serverd.AuthConfigLoader {
	return &authConfig{cfg: cfg}
}

type authConfig struct {
	cfg config.Database
}

func (ac *authConfig) LoadAuthConfig(context.Context) (*serverd.AuthConfig, error) {
	return &serverd.AuthConfig{URI: ac.cfg.URI}, nil
}
