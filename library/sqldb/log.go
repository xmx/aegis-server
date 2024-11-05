package sqldb

import (
	"log/slog"

	"github.com/go-sql-driver/mysql"
)

func NewLog(l *slog.Logger) mysql.Logger {
	return &mysqlLog{l: l}
}

type mysqlLog struct {
	l *slog.Logger
}

func (m *mysqlLog) Print(vs ...any) {
	m.l.Info("mysql", vs)
}
