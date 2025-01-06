package sqldb

import (
	"database/sql"
	"log/slog"

	"gorm.io/gorm"
)

// Open1 连接数据库。
func Open1(dsn string, driverLog *slog.Logger, opts ...gorm.Option) (*sql.DB, error) {
	//if cfg, err := msql.ParseDSN(dsn); err == nil && cfg != nil {
	//	cfg.Logger = &mysqlLog{log: driverLog}
	//	dia := &mysql.Dialector{
	//		Config: &mysql.Config{DSN: dsn, DSNConfig: cfg},
	//	}
	//	db, err := gorm.Open(dia, opts...)
	//
	//	return db, false, err
	//}
	//
	//cfg, err := pq.ParseConfig(dsn)
	//if err != nil {
	//	return nil, false, err
	//}
	//
	//cfg.Logger = &gaussLog{log: driverLog}
	//connector, err := pq.NewConnectorConfig(cfg)
	//if err != nil {
	//	return nil, true, err
	//}
	//conn := sql.OpenDB(connector)
	//gcfg := opengauss.Config{
	//	DSN:  dsn,
	//	Conn: conn,
	//}
	//dia := opengauss.New(gcfg)
	//db, err := gorm.Open(dia, opts...)
	//
	//return db, true, err

	return nil, nil
}
