package launch

import (
	"context"
	"crypto/tls"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
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
	logCfg := cfg.Logger
	logWriter := ioext.NewAttachWriter()
	if lumber := logCfg.Lumber(); lumber != nil {
		defer lumber.Close()
		logWriter.Attach(lumber)
	}
	if logCfg.Terminal {
		logWriter.Attach(os.Stdout)
	}

	logLevel := new(slog.LevelVar)
	if err := logLevel.UnmarshalText([]byte(logCfg.Level)); err != nil {
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
	parseModel(gdb)

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

	dbfs := gridfs.NewFS(qry)
	logService := service.NewLog(logLevel, logWriter, log)
	termService := service.NewTerm(log)
	fileService := service.NewFile(qry, log)

	jsLoaders := []jsvm.Loader{
		jslib.OS(),
		jslib.Time(),
		jslib.Context(),
		jslib.Console(io.Discard),
	}
	playerService := service.NewPlayer(jsLoaders, log)

	const basePath, webuiPath = "/api", "/webui"
	staticPath := path.Join(basePath, webuiPath)
	routes := []shipx.Router{
		restapi.NewConfigCertificate(configCertificateService),
		restapi.NewDAV(basePath, "/"),
		restapi.NewFile(dbfs, fileService),
		restapi.NewLog(logService),
		restapi.NewTerm(termService),
		restapi.NewPlay(playerService),
	}

	sh := ship.Default()
	sh.Validator = valid
	sh.NotFound = shipx.NotFound
	sh.HandleError = shipx.HandleError
	sh.Logger = logext.Ship(logHandler)
	sh.Route("/").GET(func(c *ship.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, staticPath)
	})
	if dir := srvCfg.Static; dir != "" {
		routes = append(routes, restapi.NewStatic(webuiPath, dir))
	}

	baseAPI := sh.Group(basePath).Use(middle.WAF(log))
	if err = shipx.BindRouters(baseAPI, routes); err != nil { // 注册路由
		return err
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
}
