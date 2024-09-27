package profile

import "time"

type Config struct {
	Database Database `json:"database"`
	Logger   Logger   `json:"logger"`
}

type Database struct {
	DSN         string   `json:"dsn"           validate:"required"`
	MaxOpenConn int      `json:"max_open_conn" validate:"gte=0"`
	MaxIdleConn int      `json:"max_idle_conn" validate:"gte=0"`
	MaxLifetime Duration `json:"max_lifetime"`
	MaxIdleTime Duration `json:"max_idle_time"`
	Migrate     bool     `json:"migrate"`
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
