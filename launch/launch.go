package launch

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"

	"github.com/quic-go/quic-go/http3"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/datalayer/query"
	"github.com/xmx/aegis-server/handler/restapi"
	"github.com/xmx/aegis-server/handler/shipx"
	"github.com/xmx/aegis-server/infra/config"
	"github.com/xmx/aegis-server/infra/gormlog"
	"github.com/xmx/aegis-server/library/credential"
	"github.com/xmx/aegis-server/library/profile"
	"github.com/xmx/aegis-server/library/sqldb"
	"github.com/xmx/aegis-server/library/validation"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Run(ctx context.Context, path string) error {
	var cfg config.Config
	if err := profile.JSON(path, &cfg); err != nil {
		return err
	}

	return Exec(ctx, cfg)
}

// Exec 运行服务。
//
//goland:noinspection GoUnhandledErrorResult
func Exec(ctx context.Context, cfg config.Config) error {
	// 创建参数校验器，并校验配置文件。
	validTags := []string{"json", "query", "form", "yaml", "xml"}
	valid := validation.NewValidator(validation.TagNameFunc(validTags))
	if err := valid.Validate(cfg); err != nil {
		return err
	}

	// 初始化日志组件。
	logOpt, logWC := cfg.Logger.Option()
	defer logWC.Close()
	logHandler := slog.NewJSONHandler(logWC, logOpt)
	log := slog.New(logHandler)

	// 连接数据库
	db, err := sqldb.TiDB(cfg.Database.DSN)
	if err != nil {
		return fmt.Errorf("连接数据库错误：%w", err)
	}
	defer db.Close()

	mysqlCfg := &mysql.Config{Conn: db}
	gormLog := gormlog.NewLog(logHandler, logger.Config{})
	gdb, err := gorm.Open(mysql.Dialector{Config: mysqlCfg}, &gorm.Config{Logger: gormLog})
	if err != nil {
		return fmt.Errorf("gorm.Open 错误：%w", err)
	}
	qry := query.Use(gdb)

	//if err = autoMigrate(gdb); err != nil {
	//	return fmt.Errorf("auto migration 错误：%w", err)
	//}

	baseTLS := &tls.Config{NextProtos: []string{"h2", "h3", "aegis"}}
	poolTLS := credential.Pool(baseTLS)

	routeRegisters := make([]shipx.Register, 0, 50)
	configCertificateService := service.ConfigCertificate(poolTLS, qry, log)
	if err = configCertificateService.Refresh(ctx); err != nil { // 初始化刷新证书池。
		return err
	}

	configCertificateAPI := restapi.ConfigCertificate(configCertificateService)
	transportAPI := restapi.Transport()
	routeRegisters = append(routeRegisters, configCertificateAPI, transportAPI)

	sh := ship.Default()
	sh.Validator = valid
	sh.NotFound = shipx.NotFound
	sh.HandleError = shipx.HandleError
	baseAPI := sh.Group("/api")
	anon := baseAPI.Clone()
	auth := baseAPI.Clone()
	route := shipx.NewRouter(anon, auth)
	for _, reg := range routeRegisters {
		if err = reg.Register(route); err != nil {
			return err
		}
	}

	srv := &http3.Server{
		Addr:    ":1443",
		Handler: sh,
		TLSConfig: &tls.Config{
			GetConfigForClient: poolTLS.Match,
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
