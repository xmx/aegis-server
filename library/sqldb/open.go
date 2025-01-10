package sqldb

import (
	"crypto/tls"
	"database/sql"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm/logger"
)

type Config struct {
	DSN         string
	MaxOpenConn int
	MaxIdleConn int
	MaxLifetime time.Duration
	MaxIdleTime time.Duration
	LogConfig   logger.Config
}

func Open(c Config, l mysql.Logger) (*sql.DB, error) {
	dsn := c.DSN
	if err := autoTLS(dsn); err != nil {
		return nil, err
	}

	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	cfg.Logger = l
	drv, err := mysql.NewConnector(cfg)
	if err != nil {
		return nil, err
	}
	db := sql.OpenDB(drv)
	db.SetMaxOpenConns(c.MaxOpenConn)
	db.SetMaxIdleConns(c.MaxIdleConn)
	db.SetConnMaxLifetime(c.MaxLifetime)
	db.SetConnMaxIdleTime(c.MaxLifetime)

	return db, nil
}

func autoTLS(dsn string) error {
	idx := strings.LastIndex(dsn, "?")
	if idx < 0 {
		return nil
	}

	params := dsn[idx+1:]
	key := lookupTLS(params)
	if key == "" {
		return nil
	}
	// FIXME 获取 dsn 中的 mysql 地址，由于 mysql 的 dsn 比较特殊，自己写一个提取器比较麻烦，
	//      复用 mysql.ParseDSN 方法来解析 dsn 并提取地址，拼接 &tls=true 参数是为了绕过 TLS 检查。
	//      该方式虽然绕过了 TLS 检查，但可能在其他特殊情况下会报错。
	//
	//      TLS 检查逻辑：https://github.com/go-sql-driver/mysql/blob/v1.8.1/dsn.go#L180-L184
	cfg, err := mysql.ParseDSN(dsn + "&tls=true")
	if err != nil {
		return err
	}

	addr := cfg.Addr
	if host, _, exx := net.SplitHostPort(addr); exx == nil {
		addr = host
	}
	tlsCfg := &tls.Config{ServerName: addr}

	return mysql.RegisterTLSConfig(key, tlsCfg)
}

// lookupTLS https://github.com/go-sql-driver/mysql/blob/v1.8.1/dsn.go#L439
//
// 与驱动处理方式保持一致：如果存在多个 tls 则取最后的。
//
// example:
//
//	lookupTLS("tls=a&tls=b&tls=c") => "c"
func lookupTLS(params string) string {
	var name string
	for _, v := range strings.Split(params, "&") {
		key, value, found := strings.Cut(v, "=")
		if !found || key != "tls" || isReservedTLS(value) {
			continue
		}
		name, _ = url.QueryUnescape(value)
	}

	return name
}

// isReservedTLS https://github.com/go-sql-driver/mysql/blob/v1.8.1/utils.go#L58-L60
func isReservedTLS(s string) bool {
	switch s {
	case "1", "true", "TRUE", "True", "0", "false", "FALSE", "False":
		return true
	default:
		lower := strings.ToLower(s)
		return lower == "skip-verify" || lower == "preferred"
	}
}
