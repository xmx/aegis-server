package launch

import (
	"context"
	"crypto/tls"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

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
	initrestapi "github.com/xmx/aegis-server/applet/initialize/restapi"
	"github.com/xmx/aegis-server/business/validext"
	"github.com/xmx/aegis-server/channel/serverd"
	"github.com/xmx/aegis-server/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/connstring"
)

func Run(ctx context.Context, cfgfile string) error {
	valid := validation.New()
	_ = valid.RegisterCustomValidations(validation.Customs())
	_ = valid.RegisterCustomValidations(validext.Customs())

	logHandlers := logger.NewHandler(logger.NewTint(os.Stdout, nil))
	log := slog.New(logHandlers)
	log.Info("程序正在启动...")

	if cfgfile == "" {
		cfgfile = config.Filename
	}
	cfr := profile.File[config.Config](cfgfile)
	if cfg, err := cfr.Read(); err == nil {
		return run(ctx, cfg, valid, logHandlers, log)
	} else if !os.IsNotExist(err) {
		log.Error("读取配置文件出错", "error", err)
		return err
	}

	log.Warn("程序等待初始化")
	results := make(chan *config.Config, 1)
	routes := []shipx.RouteRegister{
		initrestapi.NewInstall(results),
	}

	sh := ship.Default()
	sh.Validator = valid
	sh.NotFound = shipx.NotFound
	sh.HandleError = shipx.HandleError
	sh.Logger = logger.NewShip(logHandlers, 6)

	sh.Route("/").Static(config.InitialStatic)
	apiRBG := sh.Group("/api/v1")
	if err := shipx.RegisterRoutes(apiRBG, routes); err != nil {
		log.Error("注册初始化路由错误", "error", err)
		return err
	}

	// 从环境变量中获取
	const envKey = "INIT_LISTEN"
	listen := os.Getenv(envKey)
	if listen == "" {
		log.Info("如需指定监听地址，请设置环境变量", "env_key", envKey)
	}
	lis, err := net.Listen("tcp", listen)
	if err != nil {
		log.Error("初始化程序监听网络失败", "error", err)
		return err
	}
	errs := make(chan error, 1)
	srv := &http.Server{Handler: sh}
	go serveHTTP(errs, srv, lis)

	var port int
	if laddr, _ := lis.Addr().(*net.TCPAddr); laddr != nil {
		port = laddr.Port
	}
	log.Info("请打开浏览器进行初始化配置", "scheme", "http", "port", port)

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-errs:
	case cfg := <-results:
		_ = srv.Close()
		log.Warn("程序初始化完毕")
		return run(ctx, cfg, valid, logHandlers, log)
	}
	_ = srv.Close()

	return err
}

func run(ctx context.Context, cfg *config.Config, valid *validation.Validate, logh logger.Handler, log *slog.Logger) error {
	if err := valid.Validate(cfg); err != nil {
		log.Error("配置文件校验错误", "error", err)
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
		lh := logger.NewTint(os.Stdout, logOpts)
		logHandler.Attach(lh)
	}
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
	go listenHTTPS(errs, srv)
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

func serveHTTP(errs chan<- error, srv *http.Server, ln net.Listener) {
	if err := srv.Serve(ln); errors.Is(err, http.ErrServerClosed) {
		errs <- nil
	} else {
		errs <- err
	}
}

func listenQUIC(ctx context.Context, errs chan<- error, srv quick.Server) {
	errs <- srv.ListenAndServe(ctx)
}

func listenHTTPS(errs chan<- error, srv *http.Server) {
	errs <- srv.ListenAndServeTLS("", "")
}
