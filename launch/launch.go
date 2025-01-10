package launch

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/mapstruct"
	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/datalayer/gridfs"
	"github.com/xmx/aegis-server/datalayer/query"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/aegis-server/handler/middle"
	"github.com/xmx/aegis-server/handler/restapi"
	"github.com/xmx/aegis-server/handler/shipx"
	"github.com/xmx/aegis-server/library/credential"
	"github.com/xmx/aegis-server/library/multiwrite"
	"github.com/xmx/aegis-server/library/sqldb"
	"github.com/xmx/aegis-server/library/validation"
	"github.com/xmx/aegis-server/logger"
	"github.com/xmx/aegis-server/profile"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Run(ctx context.Context, path string) error {
	cfg, err := profile.JSONC(path)
	if err != nil {
		return err
	}

	return Exec(ctx, cfg)
}

// Exec 运行服务。
//
//goland:noinspection GoUnhandledErrorResult
func Exec(ctx context.Context, cfg *profile.Config) error {
	// 创建参数校验器，并校验配置文件。
	validTags := []string{"json", "query", "form", "yaml", "xml"}
	valid := validation.NewValidator(validation.TagNameFunc(validTags))
	if err := valid.Validate(cfg); err != nil {
		return err
	}

	// 初始化日志组件。
	logCfg := cfg.Logger
	logWriter := multiwrite.New(nil)
	if lumber := logCfg.Lumber(); lumber != nil {
		defer lumber.Close()
		logWriter.Attach(lumber)
	}
	if logCfg.Console {
		logWriter.Attach(os.Stdout)
	}

	logLevel := new(slog.LevelVar)
	if err := logLevel.UnmarshalText([]byte(logCfg.Level)); err != nil {
		logLevel.Set(slog.LevelWarn)
	}
	logOpt := &slog.HandlerOptions{AddSource: true, Level: logLevel}
	logHandler := slog.NewJSONHandler(logWriter, logOpt)
	log := slog.New(logHandler)

	dbCfg := mapstruct.ConfigDatabase(cfg.Database)
	gormLog, gormLevel := sqldb.NewLog(logWriter, dbCfg.LogConfig)

	// 连接数据库
	db, err := sqldb.Open(dbCfg, sqldb.NewMySQLLog(log))
	if err != nil {
		log.Error("数据库连接失败", slog.Any("error", err))
		return err
	}
	defer db.Close()

	mysqlCfg := &mysql.Config{Conn: db}
	gdb, err := gorm.Open(mysql.Dialector{Config: mysqlCfg}, &gorm.Config{Logger: gormLog})
	if err != nil {
		log.Error("数据库连接(gorm)失败", slog.Any("error", err))
		return err
	}
	qry := query.Use(gdb)

	if cfg.Database.Migrate {
		log.Info("准备检查合并数据库表结构")
		if err = autoMigrate(gdb); err != nil {
			log.Error("合并数据库错误", slog.Any("error", err))
			return err
		}
		log.Info("检查合并数据库表结构结束")
	}

	var useTLS bool
	baseTLS := &tls.Config{NextProtos: []string{"h2", "h3", "aegis"}}
	poolTLS := credential.NewPool(baseTLS)

	oplogRepo := repository.NewOplog(qry)
	oplogService := service.NewOplog(oplogRepo, log)
	configCertificateService := service.NewConfigCertificate(poolTLS, qry, log)
	if num, exx := configCertificateService.Refresh(ctx); exx != nil { // 初始化刷新证书池。
		log.Error("初始化证书错误", slog.Any("error", exx))
		return exx
	} else {
		useTLS = num > 0
	}

	dbfs := gridfs.NewFS(qry)
	logService := service.NewLog(logLevel, gormLevel, logWriter, log)
	termService := service.NewTerm(log)
	fileService := service.NewFile(qry, dbfs, log)

	const basePath = "/api"
	routes := []shipx.Router{
		restapi.NewAuth(),
		restapi.NewConfigCertificate(configCertificateService),
		restapi.NewDAV(basePath, "/"),
		restapi.NewFile(fileService),
		restapi.NewLog(logService),
		restapi.NewOplog(oplogService),
		restapi.NewTerm(termService),
		restapi.NewPlay(cfg),
	}

	srvCfg := cfg.Server
	sh := ship.Default()
	sh.Validator = valid
	sh.NotFound = shipx.NotFound
	sh.HandleError = shipx.HandleError
	sh.Logger = logger.Ship(logHandler)
	if static := srvCfg.Static; static != "" {
		sh.Route("/").Static(static)
	}

	baseAPI := sh.Group(basePath).Use(middle.WAF(oplogRepo.Create))
	if err = shipx.BindRouters(baseAPI, routes); err != nil { // 注册路由
		return err
	}

	lis, err := net.Listen("tcp", srvCfg.Addr)
	if err != nil {
		return err
	}
	if addr := listenAddr(lis); addr != "" {
		log.Warn("监听地址", slog.Any("listen", addr))
	}
	tlsCfg := &tls.Config{GetConfigForClient: poolTLS.Match}
	srv := &http.Server{Handler: sh, TLSConfig: tlsCfg}
	errs := make(chan error)
	go serveHTTP(srv, lis, useTLS, errs)
	select {
	case err = <-errs:
	case <-ctx.Done():
	}
	_ = srv.Close()
	if err != nil {
		log.Error("程序运行错误", slog.Any("error", err))
	} else {
		log.Warn("程序结束运行")
	}

	return err
}

func serveHTTP(srv *http.Server, lis net.Listener, useTLS bool, errs chan<- error) {
	if useTLS {
		errs <- srv.ServeTLS(lis, "", "")
	} else {
		errs <- srv.Serve(lis)
	}
}

func listenAddr(l net.Listener) string {
	switch v := l.(type) {
	case *net.TCPListener:
		return v.Addr().String()
	case *net.UnixListener:
		return v.Addr().String()
	default:
		return ""
	}
}
