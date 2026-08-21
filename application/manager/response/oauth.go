package response

type OAuthProvider struct {
	Provider    string   `json:"provider"`
	ClientID    string   `json:"client_id"`
	RedirectURI string   `json:"redirect_uri"`
	Scopes      []string `json:"scopes,omitzero"`
}
