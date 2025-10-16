package launch

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/robfig/cron/v3"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-common/jsos/jsmod"
	"github.com/xmx/aegis-common/jsos/jsvm"
	"github.com/xmx/aegis-common/library/cronv3"
	"github.com/xmx/aegis-common/library/httpkit"
	"github.com/xmx/aegis-common/library/validation"
	"github.com/xmx/aegis-common/logger"
	"github.com/xmx/aegis-common/profile"
	"github.com/xmx/aegis-common/shipx"
	"github.com/xmx/aegis-common/tunnel/tunutil"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-control/linkhub"
	"github.com/xmx/aegis-control/quick"
	expmiddle "github.com/xmx/aegis-server/applet/expose/middle"
	exprestapi "github.com/xmx/aegis-server/applet/expose/restapi"
	expservice "github.com/xmx/aegis-server/applet/expose/service"
	"github.com/xmx/aegis-server/business/validext"
	"github.com/xmx/aegis-server/channel/serverd"
	"github.com/xmx/aegis-server/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/connstring"
)

func Run(ctx context.Context, path string) error {
	// 2<<22 = 8388608 (8 MiB)
	opt := profile.NewOption().Limit(2 << 22).ModuleName("aegis/server/config")
	crd := profile.NewFile[config.Config](path, opt)

	return Exec(ctx, crd)
}

// Exec 运行服务。
//
//goland:noinspection GoUnhandledErrorResult
func Exec(ctx context.Context, crd profile.Reader[config.Config]) error {
	cfg, err := crd.Read(ctx)
	if err != nil {
		return err
	}

	// 创建参数校验器，并校验配置文件。
	valid := validation.New()
	valid.RegisterCustomValidations(validation.Customs())
	valid.RegisterCustomValidations(validext.Customs())
	if err = valid.Validate(cfg); err != nil {
		return err
	}

	// 初始化日志组件。
	logCfg := cfg.Logger
	logLevel := logCfg.LevelVar()
	logHandler := logger.NewHandler()
	logOpts := &slog.HandlerOptions{AddSource: true, Level: logLevel}
	if lumber := logCfg.Lumber(); lumber != nil {
		defer lumber.Close()
		logHandler.Attach(slog.NewJSONHandler(lumber, logOpts))
	}
	if logCfg.Console {
		logHandler.Attach(tint.NewHandler(os.Stdout, &tint.Options{
			AddSource:  true,
			Level:      logLevel,
			TimeFormat: time.DateTime,
		}))
	}
	log := slog.New(logHandler)
	log.Info("日志初始化完毕")

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
	log.Info("数据库连接成功")

	crond := cronv3.New(ctx, log, cron.WithSeconds())
	crond.Start()
	defer crond.Stop()

	mongoDB := cli.Database(mongoURL.Database)
	repoAll := repository.NewAll(mongoDB)

	if err = repoAll.CreateIndex(ctx); err != nil {
		return err
	}
	log.Info("数据库索引建立完毕")

	hub := linkhub.NewHub(32)
	brokerDialer := linkhub.NewSuffixDialer(hub, tunutil.BrokerHostSuffix)
	defaultDialer := tunutil.DefaultDialer()
	dialer := tunutil.NewMatchDialer(defaultDialer, brokerDialer)

	httpTrip := &http.Transport{DialContext: dialer.DialContext}
	httpCli := httpkit.NewClient(&http.Client{Transport: httpTrip})
	certificateSvc := expservice.NewCertificate(repoAll, log)
	fsSvc := expservice.NewFS(repoAll, log)
	_ = httpCli

	brokSH := ship.Default()
	brokSH.Validator = valid
	brokSH.NotFound = shipx.NotFound
	brokSH.HandleError = shipx.HandleErrorWithHost("server.internal")
	brokSH.Logger = logger.NewShip(logHandler, 6)

	{
		// aliveSvc := bservice.NewAlive(repoAll, log)
		routes := []shipx.RouteRegister{
			// brkrestapi.NewAlive(aliveSvc),
		}
		brokRGB := brokSH.Group("/api")
		if err = shipx.RegisterRoutes(brokRGB, routes); err != nil {
			return err
		}
	}

	agentSvc := expservice.NewAgent(repoAll, log)
	brokerSvc := expservice.NewBroker(repoAll, hub, log)
	if err = brokerReset(brokerSvc); err != nil {
		return err
	}

	serverdOpt := serverd.NewOption().Handler(brokSH).Logger(log).Huber(hub)
	brokerTunnelHandler := serverd.New(repoAll, cfg, serverdOpt)

	jsmodules := []jsvm.Module{
		jsmod.NewCrontab(crond),
	}
	const apiPath = "/api"
	routes := []shipx.RouteRegister{
		exprestapi.NewAgent(agentSvc),
		exprestapi.NewBroker(brokerSvc),
		exprestapi.NewCertificate(certificateSvc),
		exprestapi.NewFS(fsSvc),
		exprestapi.NewLog(logHandler),
		exprestapi.NewPlay(jsmodules),
		exprestapi.NewReverse(dialer, repoAll),
		exprestapi.NewTunnel(brokerTunnelHandler),
		exprestapi.NewDAV(apiPath, "/"),
		exprestapi.NewSystem(cfg),
		shipx.NewHealth(),
		shipx.NewPprof(),
	}

	srvCfg := cfg.Server
	outSH := ship.Default()
	outSH.Validator = valid
	outSH.NotFound = shipx.NotFound
	outSH.HandleError = shipx.HandleError
	outSH.Logger = logger.NewShip(logHandler, 6)

	rootRGB := outSH.Group("/")
	_ = exprestapi.NewStatic(srvCfg.Static).RegisterRoute(rootRGB)
	apiRGB := rootRGB.Group(apiPath).Use(expmiddle.WAF(nil))
	if err = shipx.RegisterRoutes(apiRGB, routes); err != nil { // 注册路由
		return err
	}
	log.Info("HTTP 路由注册完毕")

	listenAddr := srvCfg.Addr
	if listenAddr == "" {
		listenAddr = ":443"
	}
	tlsCfg := &tls.Config{
		GetCertificate: certificateSvc.GetCertificate,
		NextProtos:     []string{"h2", "http/1.1", "aegis"},
		MinVersion:     tls.VersionTLS13,
	}
	httpLog := logger.NewV1(slog.New(logger.Skip(logHandler, 8)))
	srv := &http.Server{
		Addr:      listenAddr,
		Handler:   outSH,
		TLSConfig: tlsCfg,
		ErrorLog:  httpLog,
	}
	quicSrv := &quick.QUICGo{
		Addr:      listenAddr,
		Handler:   brokerTunnelHandler,
		TLSConfig: tlsCfg,
	}
	log.Info("监听地址", "listen_addr", listenAddr)

	errs := make(chan error)
	go listenQUIC(ctx, errs, quicSrv)
	go listenHTTP(errs, srv)
	select {
	case err = <-errs:
	case <-ctx.Done():
	}
	_ = srv.Close()
	_ = quicSrv.Close()
	_ = brokerReset(brokerSvc)

	if err != nil {
		log.Error("程序运行错误", slog.Any("error", err))
	} else {
		log.Warn("程序结束运行")
	}

	return err
}

func disconnectDB(cli *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = cli.Disconnect(ctx)
}

func brokerReset(brk *expservice.Broker) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return brk.Reset(ctx)
}

func listenHTTP(errs chan<- error, srv *http.Server) {
	errs <- srv.ListenAndServeTLS("", "")
}

func listenQUIC(ctx context.Context, errs chan<- error, srv quick.Server) {
	errs <- srv.ListenAndServe(ctx)
}
