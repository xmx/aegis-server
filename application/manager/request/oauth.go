package request

type OAuthCode struct {
	Code string `json:"code" validate:"required,lte=500"`
}

type OAuthProvider struct {
	Provider string `json:"provider" query:"provider" validate:"required,lte=255"`
	Origin   string `json:"origin"   query:"origin"   validate:"lte=255"`
}
