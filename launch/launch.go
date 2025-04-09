package launch

import (
	"context"
	"crypto/tls"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/xmx/aegis-server/business/jsext"

	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/aegis-server/handler/middle"
	"github.com/xmx/aegis-server/handler/restapi"
	"github.com/xmx/aegis-server/handler/shipx"
	"github.com/xmx/aegis-server/jsrun/jsmod"
	"github.com/xmx/aegis-server/jsrun/jsvm"
	"github.com/xmx/aegis-server/library/credential"
	"github.com/xmx/aegis-server/library/cronv3"
	"github.com/xmx/aegis-server/library/dynwriter"
	"github.com/xmx/aegis-server/library/validation"
	"github.com/xmx/aegis-server/logger"
	"github.com/xmx/aegis-server/profile"
	"github.com/xmx/ship"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/connstring"
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
	valid := validation.New(validation.TagNameFunc([]string{"json"}))
	if err := valid.Validate(ctx, cfg); err != nil {
		return err
	}

	// 初始化日志组件。
	logCfg := cfg.Logger
	logWriter := dynwriter.New()
	if lumber := logCfg.Lumber(); lumber != nil {
		defer lumber.Close()
		logWriter.Attach(lumber)
	}
	if logCfg.Console {
		logWriter.Attach(os.Stdout)
	}

	logLevel := new(slog.LevelVar)
	if err := logLevel.UnmarshalText([]byte(logCfg.Level)); err != nil {
		logLevel.Set(slog.LevelInfo)
	}
	logOpt := &slog.HandlerOptions{AddSource: true, Level: logLevel}
	logHandler := slog.NewJSONHandler(logWriter, logOpt)
	log := slog.New(logHandler)

	// -----[ 初始化 mongodb ]-----
	mongoURI := cfg.Database.URI
	mongoURL, err := connstring.ParseAndValidate(mongoURI)
	if err != nil {
		return err
	}

	mongoLogOpt := options.Logger().
		SetSink(logger.NewSink(logHandler, 13)).
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)
	mongoOpt := options.Client().
		ApplyURI(mongoURI).
		SetLoggerOptions(mongoLogOpt)
	cli, err := mongo.Connect(mongoOpt)
	if err != nil {
		return err
	}
	defer disconnectDB(cli)

	cronLog := slog.New(logger.Skip(logHandler, 5))
	crontab := cronv3.New(cronLog)
	crontab.Start()
	defer crontab.Stop()

	mongoDB := cli.Database(mongoURL.Database)
	repoAll := repository.NewAll(mongoDB)
	if err = repoAll.CreateIndex(ctx); err != nil {
		return err
	}
	registerValidator(valid, repoAll, log)

	var useTLS bool
	baseTLS := &tls.Config{NextProtos: []string{"h2", "h3", "aegis"}}
	poolTLS := credential.NewPool(baseTLS)

	certificateSvc, err := service.NewCertificate(repoAll, poolTLS, log)
	if err != nil {
		return err
	}
	logSvc := service.NewLog(logLevel, logWriter, log)

	if num, exx := certificateSvc.Refresh(ctx); exx != nil { // 初始化刷新证书池。
		log.Error("初始化证书错误", slog.Any("error", exx))
		return exx
	} else {
		useTLS = num > 0
	}
	termSvc := service.NewTerm(log)

	const basePath = "/api"
	modules := []jsvm.GlobalRegister{
		logSvc,
		jsmod.NewConsole(io.Discard, io.Discard),
		jsmod.NewContext(),
		jsmod.NewExec(),
		jsmod.NewIO(),
		jsmod.NewOS(),
		jsmod.NewNet(),
		jsmod.NewRuntime(),
		jsmod.NewTime(),
		jsext.NewCrontab(crontab),
	}
	routes := []shipx.RouteRegister{
		restapi.NewAuth(),
		restapi.NewCertificate(certificateSvc),
		restapi.NewLog(logSvc),
		restapi.NewDAV(basePath, "/"),
		restapi.NewPlay(modules),
		restapi.NewSystem(),
		restapi.NewTerm(termSvc),
		// restapi.NewFile(fileSvc),
		// restapi.NewLog(logSvc),
		// restapi.NewOplog(oplogSvc),
		// restapi.NewTerm(termSvc),
		// restapi.NewPlay(modules),
	}

	srvCfg := cfg.Server
	sh := ship.Default()
	sh.Validator = valid
	sh.NotFound = shipx.NotFound
	sh.HandleError = shipx.HandleError
	sh.Logger = logger.NewShip(logHandler, 6)
	if static := srvCfg.Static; static != "" {
		sh.Route("/").Static(static)
	}

	baseAPI := sh.Group(basePath).Use(middle.WAF(nil))
	if err = shipx.RegisterRoutes(baseAPI, routes); err != nil { // 注册路由
		return err
	}

	lc := new(net.ListenConfig)
	lc.SetMultipathTCP(true)
	lis, err := lc.Listen(ctx, "tcp", srvCfg.Addr)
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

func disconnectDB(cli *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = cli.Disconnect(ctx)
}
