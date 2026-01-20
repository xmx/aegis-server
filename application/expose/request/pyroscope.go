package request

type PyroscopeUpsert struct {
	Name     string `json:"name"     validate:"required"`
	Address  string `json:"address"  validate:"required,http_url"`
	Username string `json:"username" validate:"lte=20"`
	Password string `json:"password" validate:"lte=100"`
	Enabled  bool   `json:"enabled"`
}
