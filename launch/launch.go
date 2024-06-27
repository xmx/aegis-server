package launch

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"os"
	"time"

	"github.com/quic-go/quic-go/http3"
	"github.com/xmx/aegis-server/datalayer/query"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/aegis-server/library/sqldb"
	"github.com/xmx/aegis-server/memconf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Run(ctx context.Context, path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var jsonCfg struct {
		DSN string `json:"dsn"`
	}
	if err = json.Unmarshal(file, &jsonCfg); err != nil {
		return err
	}

	return Exec(ctx, jsonCfg.DSN)
}

// Exec 运行服务。
//
//goland:noinspection GoUnhandledErrorResult
func Exec(ctx context.Context, dsn string) error {
	db, err := sqldb.TiDB(dsn, 10*time.Second)
	if err != nil {
		return err
	}
	defer db.Close()

	mysqlCfg := &mysql.Config{Conn: db}
	gdb, err := gorm.Open(mysql.Dialector{Config: mysqlCfg})
	if err != nil {
		return err
	}

	if err = autoMigrate(gdb); err != nil {
		return err
	}
	qry := query.Use(gdb)

	configCertificateRepository := repository.ConfigCertificate(qry)
	configCertificateConfigurer := memconf.ConfigCertificate(configCertificateRepository)

	srv := &http3.Server{
		Addr: ":1443",
		TLSConfig: &tls.Config{
			GetConfigForClient: configCertificateConfigurer.Certificate,
		},
	}
	errs := make(chan error)
	go serveHTTP(srv, errs)

	select {
	case err = <-errs:
	case <-ctx.Done():
	}
	_ = srv.Close()

	return err
}

func serveHTTP(srv *http3.Server, errs chan<- error) {
	errs <- srv.ListenAndServe()
}
