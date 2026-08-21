package response

import "strconv"

type OAuthProvider struct {
	Provider    string   `json:"provider"`
	ClientID    string   `json:"client_id"`
	RedirectURI string   `json:"redirect_uri"`
	Scopes      []string `json:"scopes,omitzero"`
	AuthURL     string   `json:"auth_url"`
}

type OAuthGitHub struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name,omitzero"`
	AvatarURL string `json:"avatar_url,omitzero"`
	Company   string `json:"company,omitzero"`
	Email     string `json:"email,omitzero"`
	Location  string `json:"location,omitzero"`
}

func (o OAuthGitHub) PUID() string {
	return strconv.Itoa(o.ID)
}
