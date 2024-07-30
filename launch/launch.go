package launch

import (
	"context"
	"crypto/tls"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/xmx/aegis-server/jsenv/jslib"
	"github.com/xmx/aegis-server/jsenv/jsvm"

	"github.com/xgfone/ship/v5"
	"github.com/xgfone/ship/v5/middleware"
	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/datalayer/query"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/aegis-server/handler/middle"
	"github.com/xmx/aegis-server/handler/restapi"
	"github.com/xmx/aegis-server/handler/shipx"
	"github.com/xmx/aegis-server/infra/logext"
	"github.com/xmx/aegis-server/infra/profile"
	"github.com/xmx/aegis-server/library/credential"
	"github.com/xmx/aegis-server/library/sqldb"
	"github.com/xmx/aegis-server/library/validation"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Run(ctx context.Context, path string) error {
	cfg, err := profile.JSON(path)
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
	logWC := cfg.Logger.Writer()
	defer logWC.Close()

	logLevel := new(slog.LevelVar)
	logLevel.Set(slog.LevelWarn) // 默认日志级别
	_ = logLevel.UnmarshalText([]byte(cfg.Logger.Level))
	logOpt := &slog.HandlerOptions{AddSource: true, Level: logLevel}

	logHandler := slog.NewJSONHandler(logWC, logOpt)
	log := slog.New(logHandler)

	// 连接数据库
	db, err := sqldb.TiDB(cfg.Database.TiDB())
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
	configServerRepository := repository.ConfigServer(qry)
	srvCfg, err := configServerRepository.Enabled(ctx)
	if err != nil {
		log.Error("查询 server 配置错误", slog.Any("error", err))
		return err
	}

	baseTLS := &tls.Config{NextProtos: []string{"h2", "h3", "aegis"}}
	poolTLS := credential.Pool(baseTLS)

	routeRegisters := make([]shipx.Register, 0, 50)
	configCertificateService := service.ConfigCertificate(poolTLS, qry, log)
	if err = configCertificateService.Refresh(ctx); err != nil { // 初始化刷新证书池。
		log.Error("初始化证书错误", slog.Any("error", err))
		return err
	}

	configCertificateAPI := restapi.ConfigCertificate(configCertificateService)
	logAPI := restapi.Log(logWC, logLevel)
	routeRegisters = append(routeRegisters, configCertificateAPI, logAPI)

	{
		loads := []jsvm.Loader{
			jslib.OS(),
			jslib.Time(),
			jslib.Context(),
			jslib.Console(io.Discard), // 认识丢弃输出数据。
		}
		playAPI := restapi.Play(loads)
		routeRegisters = append(routeRegisters, playAPI)
	}

	sh := ship.Default()
	sh.Validator = valid
	sh.NotFound = shipx.NotFound
	sh.HandleError = shipx.HandleError
	sh.Logger = logext.Ship(logHandler)
	if dir := srvCfg.Static; dir != "" {
		sh.Route("/").Static(dir)
	}

	baseAPI := sh.Group("/api").Use(middle.WAF(log), middleware.CORS(nil))
	anon := baseAPI.Clone()
	auth := baseAPI.Clone()

	for _, reg := range routeRegisters {
		route := shipx.NewRouter(anon, auth)
		if err = reg.Register(route); err != nil {
			log.Error("注册路由错误", slog.Any("error", err))
			return err
		}
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

	return err
}

func serveHTTP(srv *http.Server, errs chan<- error) {
	errs <- srv.ListenAndServeTLS("", "")
}
