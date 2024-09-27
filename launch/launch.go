package launch

import (
	"context"
	"crypto/tls"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/argument/mapstruct"
	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/datalayer/gridfs"
	"github.com/xmx/aegis-server/datalayer/query"
	"github.com/xmx/aegis-server/handler/middle"
	"github.com/xmx/aegis-server/handler/restapi"
	"github.com/xmx/aegis-server/handler/shipx"
	"github.com/xmx/aegis-server/infra/logext"
	"github.com/xmx/aegis-server/infra/profile"
	"github.com/xmx/aegis-server/jsenv/jslib"
	"github.com/xmx/aegis-server/jsenv/jsvm"
	"github.com/xmx/aegis-server/library/credential"
	"github.com/xmx/aegis-server/library/ioext"
	"github.com/xmx/aegis-server/library/sqldb"
	"github.com/xmx/aegis-server/library/validation"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	logWriter := ioext.NewAttachWriter()
	if lumber := cfg.Logger.Lumber(); lumber != nil {
		defer lumber.Close()
		logWriter.Attach(lumber)
	}
	if cfg.Logger.Terminal {
		logWriter.Attach(os.Stdout)
	}

	logLevel := new(slog.LevelVar)
	if err := logLevel.UnmarshalText([]byte(cfg.Logger.Level)); err != nil {
		logLevel.Set(slog.LevelWarn)
	}
	logOpt := &slog.HandlerOptions{AddSource: true, Level: logLevel}
	logHandler := slog.NewJSONHandler(logWriter, logOpt)
	log := slog.New(logHandler)

	// 连接数据库
	dbCfg := mapstruct.ConfigDatabase(cfg.Database)
	db, err := sqldb.TiDB(dbCfg)
	if err != nil {
		log.Error("数据库连接失败", slog.Any("error", err))
		return err
	}
	defer db.Close()

	glogCfg := logger.Config{SlowThreshold: 300 * time.Millisecond, LogLevel: logger.Info}
	gormLog := logext.Gorm(logHandler, glogCfg)
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

	// 查询 server 配置
	configServerService := service.NewConfigServer(qry)
	srvCfg, err := configServerService.Enabled(ctx)
	if err != nil {
		log.Error("查询 server 配置错误", slog.Any("error", err))
		return err
	}

	baseTLS := &tls.Config{NextProtos: []string{"h2", "h3", "aegis"}}
	poolTLS := credential.Pool(baseTLS)

	configCertificateService := service.NewConfigCertificate(poolTLS, qry, log)
	if err = configCertificateService.Refresh(ctx); err != nil { // 初始化刷新证书池。
		log.Error("初始化证书错误", slog.Any("error", err))
		return err
	}

	name := qry.GridFile.ID.ColumnName()
	qry.GridFile.ID.Desc()
	log.Info(string(name))
	log.Info(name.String())

	dbfs := gridfs.NewFS(qry)

	logService := service.NewLog(logLevel, logWriter, log)
	termService := service.NewTerm(log)
	fileService := service.NewFile(qry, log)

	configCertificateAPI := restapi.NewConfigCertificate(configCertificateService)
	logAPI := restapi.NewLog(logService)
	termAPI := restapi.NewTerm(termService)
	vncAPI := restapi.NewVNC()
	fileAPI := restapi.NewFile(dbfs, fileService)

	routeRegisters := make([]shipx.Controller, 0, 50)
	routeRegisters = append(routeRegisters, configCertificateAPI, logAPI, termAPI, vncAPI, fileAPI)

	{
		loads := []jsvm.Loader{
			jslib.OS(),
			jslib.Time(),
			jslib.Context(),
			jslib.Console(io.Discard), // 默认丢弃输出数据。
		}
		playerService := service.NewPlayer(loads, log)
		playAPI := restapi.NewPlay(playerService)
		routeRegisters = append(routeRegisters, playAPI)
	}

	sh := ship.Default()
	sh.Validator = valid
	sh.NotFound = shipx.NotFound
	sh.HandleError = shipx.HandleError
	sh.Logger = logext.Ship(logHandler)
	sh.Route("/").GET(func(c *ship.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/api/webui/")
	})
	if dir := srvCfg.Static; dir != "" {
		staticAPI := restapi.NewStatic(dir)
		routeRegisters = append(routeRegisters, staticAPI)
	}

	const basePrefix = "/api"
	davAPI := restapi.NewDAV(basePrefix, "/", false)
	routeRegisters = append(routeRegisters, davAPI)

	baseAPI := sh.Group(basePrefix).Use(middle.WAF(log))
	anon := baseAPI.Clone()
	auth := baseAPI.Clone()

	for _, reg := range routeRegisters {
		route := shipx.NewRouter(anon, auth)
		if err = reg.Register(route); err != nil {
			log.Error("注册路由错误", slog.Any("error", err))
			return err
		}
	}
	for _, route := range sh.Routes() {
		log.Info("路由", slog.Any("route", route))
	}

	srv := &http.Server{
		Addr:      srvCfg.Addr,
		Handler:   sh,
		TLSConfig: &tls.Config{GetConfigForClient: poolTLS.Match},
	}
	errs := make(chan error)
	go serveHTTP(srv, errs)
	// go save(qry)
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

func serveHTTP(srv *http.Server, errs chan<- error) {
	errs <- srv.ListenAndServeTLS("", "")
	// errs <- srv.ListenAndServe()
}
