package launch

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/quic-go/quic-go"

	"github.com/quic-go/quic-go/http3"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/datalayer/query"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/aegis-server/handler/restapi"
	"github.com/xmx/aegis-server/handler/shipx"
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
	db, err := sqldb.TiDB(dsn)
	if err != nil {
		return fmt.Errorf("连接数据库错误：%w", err)
	}
	defer db.Close()

	mysqlCfg := &mysql.Config{Conn: db}
	gdb, err := gorm.Open(mysql.Dialector{Config: mysqlCfg})
	if err != nil {
		return fmt.Errorf("gorm.Open 错误：%w", err)
	}

	if err = autoMigrate(gdb); err != nil {
		return fmt.Errorf("auto migration 错误：%w", err)
	}
	qry := query.Use(gdb)

	//configLoggerRepository := repository.ConfigLogger(qry)
	//loggerConfig, err := configLoggerRepository.Enabled(ctx)
	//if err != nil {
	//	return err
	//}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	log := slog.New(handler)

	routeRegisters := make([]shipx.Register, 0, 50)
	configCertificateRepository := repository.ConfigCertificate(qry)
	configCertificateConfigurer := memconf.ConfigCertificate(configCertificateRepository)
	configCertificateService := service.ConfigCertificate(qry, configCertificateConfigurer, log)
	configCertificateAPI := restapi.ConfigCertificate(configCertificateService)
	routeRegisters = append(routeRegisters, configCertificateAPI)

	sh := ship.Default()
	anon := sh.Group("/").Clone()
	auth := sh.Group("/").Clone()
	route := shipx.NewRouter(anon, auth)
	for _, reg := range routeRegisters {
		if err = reg.Register(route); err != nil {
			return err
		}
	}

	srv := &http3.Server{
		Handler: sh,
	}
	errs := make(chan error)
	go serveHTTP(srv, errs)
	srv.TLSConfig.NextProtos = append(srv.TLSConfig.NextProtos, "hi")

	tlsCfg := &tls.Config{
		NextProtos:     []string{"hi"},
		GetCertificate: configCertificateConfigurer.Certificate,
	}
	listener, err := quic.ListenAddrEarly(":1443", tlsCfg, nil)
	if err != nil {
		return err
	}

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

func listenQUIC(ctx context.Context, lis *quic.EarlyListener, errs chan<- error) {
	for {
		conn, err := lis.Accept(ctx)
		if err != nil {
			errs <- err
			break
		}

	}
}
