package mapstruct

import (
	"strings"
	"time"

	"github.com/xmx/aegis-server/library/sqldb"
	"github.com/xmx/aegis-server/profile"
	"gorm.io/gorm/logger"
)

func ConfigDatabase(c profile.Database) sqldb.Config {
	logLevel := logger.Warn
	switch strings.ToUpper(c.LogLevel) {
	case "INFO":
		logLevel = logger.Info
	case "ERROR":
		logLevel = logger.Error
	}

	return sqldb.Config{
		DSN:         c.DSN,
		MaxOpenConn: c.MaxOpenConn,
		MaxIdleConn: c.MaxIdleConn,
		MaxLifetime: time.Duration(c.MaxLifetime),
		MaxIdleTime: time.Duration(c.MaxIdleTime),
		LogConfig: logger.Config{
			SlowThreshold:             time.Duration(c.SlowSQL),
			IgnoreRecordNotFoundError: c.IgnoreRecordNotFoundError,
			ParameterizedQueries:      c.ParameterizedQueries,
			LogLevel:                  logLevel,
		},
	}
}
