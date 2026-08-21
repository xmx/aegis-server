package message

import "github.com/xmx/aegis-server/application/config"

type OOBEData struct {
	Config *config.Config `json:"config" validate:"required"`
	OAuth  OOBEOAuth      `json:"oauth"` // OAuth 认证配置
	Admin  OOBEAdmin      `json:"admin"`
}

type OOBEAdmin struct {
	JobNumber string `json:"job_number" validate:"required,job_number"` // 管理员工号
	Name      string `json:"name"       validate:"lte=10"`              // 管理员名字
}

type OOBEOAuth struct {
	URL         string `json:"url"          validate:"http_url"` // OAuth 认证地址
	ClientID    string `json:"client_id"    validate:"required"` // OAuth ClientID
	Secret      string `json:"secret"       validate:"required"` // OAuth ClientSecret
	RedirectURL string `json:"redirect_url" validate:"http_url"` // OAuth RedirectURL
}
