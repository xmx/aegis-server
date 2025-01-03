package mapstruct

import (
	"time"

	"github.com/xmx/aegis-server/library/sqldb"
	"github.com/xmx/aegis-server/profile"
)

func ConfigDatabase(c profile.Database) sqldb.Config {
	return sqldb.Config{
		DSN:         c.DSN,
		MaxOpenConn: c.MaxOpenConn,
		MaxIdleConn: c.MaxIdleConn,
		MaxLifetime: time.Duration(c.MaxLifetime),
		MaxIdleTime: time.Duration(c.MaxIdleTime),
	}
}
