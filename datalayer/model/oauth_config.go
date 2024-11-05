package model

type OAuthConfig struct {
	ID           int64  `json:"id,string"     gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	Name         string `json:"name"          gorm:"column:name;type:varchar(20);comment:名字"`
	Endpoint     string `json:"endpoint"      gorm:"column:endpoint;type:varchar(100);comment:Endpoint"`
	ClientID     string `json:"client_id"     gorm:"column:client_id;type:varchar(100);comment:Client ID"`
	ClientSecret string `json:"client_secret" gorm:"column:client_secret;type:varchar(255);comment:Client Secret"`
	RedirectURI  string `json:"redirect_uri"  gorm:"column:redirect_uri;type:varchar(255);comment:Redirect URI"`
}

func (OAuthConfig) TableName() string {
	return "oauth_config"
}
