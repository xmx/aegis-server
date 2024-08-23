package profile

import (
	"time"

	"github.com/xmx/aegis-server/library/sqldb"
)

type Config struct {
	Server   Server   `json:"server"`
	Database Database `json:"database"`
	Logger   Logger   `json:"logger"`
}

type Server struct {
	Addr string `json:"addr"`
}

type Database struct {
	DSN         string   `json:"dsn"           validate:"required"`
	MaxOpenConn int      `json:"max_open_conn"`
	MaxIdleConn int      `json:"max_idle_conn"`
	MaxLifetime Duration `json:"max_lifetime"`
	MaxIdleTime Duration `json:"max_idle_time"`
	Migrate     bool     `json:"migrate"`
}

func (d Database) TiDB() sqldb.Config {
	return sqldb.Config{
		DSN:         d.DSN,
		MaxOpenConn: d.MaxOpenConn,
		MaxIdleConn: d.MaxIdleConn,
		MaxLifetime: time.Duration(d.MaxLifetime),
		MaxIdleTime: time.Duration(d.MaxIdleTime),
	}
}

type Duration time.Duration

func (d *Duration) UnmarshalText(raw []byte) error {
	du, err := time.ParseDuration(string(raw))
	if err != nil {
		return err
	}
	*d = Duration(du)

	return nil
}
