package profile

import "time"

type Config struct {
	Active   string   `json:"active"`
	Server   Server   `json:"server"`
	Database Database `json:"database"`
	Logger   Logger   `json:"logger"`
}

type Database struct {
	DSN                       string   `json:"dsn"           validate:"required"`
	MaxOpenConn               int      `json:"max_open_conn" validate:"gte=0"`
	MaxIdleConn               int      `json:"max_idle_conn" validate:"gte=0"`
	MaxLifetime               Duration `json:"max_lifetime"`
	MaxIdleTime               Duration `json:"max_idle_time"`
	Migrate                   bool     `json:"migrate"`
	SlowSQL                   Duration `json:"slow_sql"` // 慢 SQL 阈值
	IgnoreRecordNotFoundError bool     `json:"ignore_record_not_found_error"`
	ParameterizedQueries      bool     `json:"parameterized_queries"`
	LogLevel                  string   `json:"log_level"`
}

type Duration time.Duration

func (du *Duration) UnmarshalText(raw []byte) error {
	d, err := time.ParseDuration(string(raw))
	if err != nil {
		return err
	}
	*du = Duration(d)

	return nil
}

type Server struct {
	Addr   string `json:"addr"   validate:"lte=100"`
	Cert   string `json:"cert"   validate:"lte=255"`
	Pkey   string `json:"pkey"   validate:"lte=255"`
	Static string `json:"static" validate:"lte=255"`
}
