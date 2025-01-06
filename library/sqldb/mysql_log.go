package sqldb

import (
	"log/slog"

	"github.com/go-sql-driver/mysql"
)

func NewMySQLLog(l *slog.Logger) mysql.Logger {
	return &mysqlLog{l: l}
}

type mysqlLog struct {
	l *slog.Logger
}

func (m *mysqlLog) Print(vs ...any) {
	size := len(vs)
	if size == 0 {
		return
	}

	msg := "mysql"
	arg0 := vs[0]
	switch v := arg0.(type) {
	case error:
		msg = v.Error()
		vs = vs[1:]
	case string:
		msg = v
		vs = vs[1:]
	}
	if len(vs) == 0 {
		m.l.Info(msg)
	} else {
		m.l.Info(msg, slog.Any("data", vs))
	}
}
