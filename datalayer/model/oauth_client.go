package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type OAuthClient struct {
	ID           bson.ObjectID `bson:"_id,omitempty"        json:"-"`
	Provider     string        `bson:"provider"             json:"provider"`
	Enabled      bool          `bson:"enabled"              json:"enabled"`
	ClientID     string        `bson:"client_id"            json:"client_id"`
	ClientSecret string        `bson:"client_secret"        json:"client_secret"`
	RedirectURIs []string      `bson:"redirect_uris"        json:"redirect_uris"`
	Scopes       []string      `bson:"scopes,omitempty"     json:"scopes,omitzero"`
	CreatedAt    time.Time     `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at,omitempty" json:"updated_at"`
}

func (OAuthClient) CollectionInfo() (string, string) { return "oauth_client", "OAuth 客户端凭证" }
