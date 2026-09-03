package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/xmx/aegis-server/application/errcode"
	"github.com/xmx/aegis-server/application/manager/response"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type OAuth struct {
	db  *repository.BaseDB
	cli *http.Client
	log *slog.Logger
}

func NewOAuth(db *repository.BaseDB, cli *http.Client, log *slog.Logger) *OAuth {
	return &OAuth{
		db:  db,
		log: log,
	}
}

func (ath *OAuth) Provider(ctx context.Context, provider string) (*model.OAuthClient, error) {
	coll := ath.db.OAuthClient()
	dat, err := coll.FindByProvider(ctx, provider)
	if err != nil {
		return nil, errcode.ErrUnsupportedOAuthProvider
	} else if !dat.Enabled {
		return nil, errcode.ErrForbiddenOAuthProvider
	}

	return dat, nil
}

func (ath *OAuth) GitHub(ctx context.Context, code string, redirectURI string) (any, error) {
	dat, err := ath.Provider(ctx, model.OAuthGithub)
	if err != nil {
		return nil, err
	}
	cfg := &oauth2.Config{
		ClientID:     dat.ClientID,
		ClientSecret: dat.ClientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  redirectURI,
		Scopes:       dat.Scopes,
	}
	info, err := ath.githubUserInfo(ctx, cfg, code)
	if err != nil {
		return nil, err
	}

	// 新增或创建
	now := time.Now()
	usr := &model.User{
		Enabled:   true,
		Provider:  model.OAuthGithub,
		PUID:      info.PUID(),
		Login:     info.Login,
		Name:      info.Name,
		AvatarURL: info.AvatarURL,
		Company:   info.Company,
		Email:     info.Email,
		Location:  info.Location,
		UpdatedAt: now,
	}

	filter := bson.M{"provider": model.OAuthGithub, "puid": usr.PUID}
	update := bson.M{"$set": usr, "$setOnInsert": bson.M{"created_at": now}}
	opts := options.UpdateOne().SetUpsert(true)
	coll := ath.db.User()
	if _, err = coll.UpdateOne(ctx, filter, update, opts); err != nil {
		return nil, err
	}

	return usr, nil
}

func (ath *OAuth) githubUserInfo(parent context.Context, cfg *oauth2.Config, code string) (*response.OAuthGitHub, error) {
	tok, err := cfg.Exchange(parent, code)
	if err != nil {
		ath.log.ErrorContext(parent, "GitHub 密钥交换出错", "err", err)
		return nil, err
	}

	ctx := parent
	const infoURL = "https://api.github.com/user"
	if ath.cli != nil {
		ctx = context.WithValue(parent, oauth2.HTTPClient, ath.cli)
	}
	cli := cfg.Client(ctx, tok)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, infoURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	info := new(response.OAuthGitHub)
	if err = json.NewDecoder(res.Body).Decode(info); err != nil {
		return nil, err
	}

	return info, nil
}
