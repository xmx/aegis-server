package launch

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-common/jsos/jsvm"
	"github.com/xmx/aegis-common/library/cronv3"
	"github.com/xmx/aegis-common/library/httpkit"
	"github.com/xmx/aegis-common/library/validation"
	"github.com/xmx/aegis-common/logger"
	"github.com/xmx/aegis-common/muxlink/muxproto"
	"github.com/xmx/aegis-common/profile"
	"github.com/xmx/aegis-common/shipx"
	"github.com/xmx/aegis-control/datalayer/repository"
	"github.com/xmx/aegis-control/linkhub"
	"github.com/xmx/aegis-control/mongodb"
	"github.com/xmx/aegis-control/quick"
	"github.com/xmx/aegis-control/tlscert"
	brkrestapi "github.com/xmx/aegis-server/application/broker/restapi"
	brkservice "github.com/xmx/aegis-server/application/broker/service"
	"github.com/xmx/aegis-server/application/crontab"
	expmiddle "github.com/xmx/aegis-server/application/expose/middle"
	exprestapi "github.com/xmx/aegis-server/application/expose/restapi"
	expservice "github.com/xmx/aegis-server/application/expose/service"
	"github.com/xmx/aegis-server/application/serverd"
	"github.com/xmx/aegis-server/application/tunutil"
	"github.com/xmx/aegis-server/application/validext"
	wizrestapi "github.com/xmx/aegis-server/application/wizard/restapi"
	wizstatic "github.com/xmx/aegis-server/application/wizard/static"
	"github.com/xmx/aegis-server/config"
	"github.com/xmx/aegis-server/library/firewalld"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/net/quic"
)

func Run(ctx context.Context, cfgFile string) error {
	valid := validation.New()
	_ = valid.RegisterCustomValidations(validation.All())
	_ = valid.RegisterCustomValidations(validext.All())

	logOpts := &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}
	logh := logger.NewMultiHandler(logger.NewTint(os.Stdout, logOpts))
	log := slog.New(logh)
	log.Info("程序正在启动...")

	cfr := profile.File[config.Config](cfgFile)
	if cfg, err := cfr.Read(); err == nil {
		return run(ctx, cfg, valid, logh, log)
	} else if !os.IsNotExist(err) {
		log.Error("读取配置文件出错", "error", err)
		return err
	}

	log.Warn("程序等待初始化")
	results := make(chan *config.Config, 1)
	routes := []shipx.RouteRegister{
		wizrestapi.NewInstall(cfgFile, results),
	}

	sh := ship.Default()
	sh.Validator = valid
	sh.NotFound = shipx.NotFound
	sh.HandleError = shipx.HandleError
	sh.Logger = logger.NewShip(logh)

	sh.Route("/").StaticFS(http.FS(wizstatic.FS))
	apiRBG := sh.Group("/api")
	if err := shipx.RegisterRoutes(apiRBG, routes); err != nil {
		log.Error("注册初始化路由错误", "error", err)
		return err
	}

	const envKey = "AEGIS_INIT_ADDR"
	listen := os.Getenv(envKey)
	if listen == "" {
		// 如果没有指定初始化监听地址，就默认监听 443，如果端口冲突或受限的网络中，
		// 需要指定特定的端口，请使用环境变量指明。
		listen = ":443"
		log.Info("如需指定监听地址，请设置环境变量", "env_key", envKey)
	}

	noneTLS := func(context.Context) ([]*tls.Certificate, error) { return nil, nil }
	tempTLS := tlscert.NewCertPool(noneTLS, true, log)
	srv := &http.Server{
		Addr:           listen,
		Handler:        sh,
		MaxHeaderBytes: 10 * 1024,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		TLSConfig:      &tls.Config{GetCertificate: tempTLS.Match},
	}

	log.Info("请打开浏览器进行初始化配置", "scheme", "https", "listen", listen)
	errs := make(chan error, 1)
	go listenHTTPS(errs, srv)

	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-errs:
	case cfg := <-results:
		cctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		_ = srv.Shutdown(cctx)
		cancel()
		log.Warn("程序初始化完毕")
		return run(ctx, cfg, valid, logh, log)
	}

	return err
}

//goland:noinspection GoUnhandledErrorResult
func run(ctx context.Context, cfg *config.Config, valid *validation.Validate, logh logger.Handler, log *slog.Logger) error {
	if err := valid.Validate(cfg); err != nil {
		log.Error("配置文件校验错误", "error", err)
		return err
	}

	// 初始化日志组件。
	logCfg := cfg.Logger
	logLevel := logCfg.LevelVar()
	logOpts := &slog.HandlerOptions{AddSource: true, Level: logLevel}
	logh.Replace() // reset 日志
	if lumber := logCfg.Lumber(); lumber != nil {
		defer lumber.Close()
		logh.Attach(slog.NewJSONHandler(lumber, logOpts))
	}
	if logCfg.Console {
		lh := logger.NewTint(os.Stdout, logOpts)
		logh.Attach(lh)
	}
	log.Info("日志初始化完毕")

	mongoLogOpt := options.Logger().
		SetSink(logger.NewSink(logh, 13)).
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)
	mongoOpt := options.Client().
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
		SetLoggerOptions(mongoLogOpt)
	db, err := mongodb.Open(cfg.Database.URI, mongoOpt)
	if err != nil {
		log.Error("数据库连接错误", "error", err)
		return err
	}
	defer disconnectDB(db)
	log.Info("数据库连接成功")

	log.Info("开始初始化数据库索引...")
	repoAll := repository.NewAll(db, log)
	if err = repoAll.CreateIndex(ctx); err != nil {
		log.Error("索引创建错误", "error", err)
		return err
	}
	log.Info("数据库索引建立完毕")

	crond := cronv3.New(ctx, log, cron.WithSeconds())
	crond.Start()
	defer crond.Stop()

	hub := linkhub.NewHub(32)
	netDialer := &net.Dialer{Timeout: 30 * time.Second}
	tunDialer := []muxproto.Dialer{
		linkhub.NewSuffixDialer(muxproto.BrokerHostSuffix, hub),
		serverd.NewFindAgentDialer(muxproto.AgentHostSuffix, hub, repoAll),
	}
	dualDialer := muxproto.NewFirstMatchDialer(tunDialer, netDialer)
	loadCert := repoAll.Certificate().Enables
	certPool := tlscert.NewCertPool(loadCert, true, log)

	httpTrip := &http.Transport{DialContext: dualDialer.DialContext}
	httpClient := &http.Client{Transport: httpTrip}
	_ = httpkit.NewClient(httpClient)
	certificateSvc := expservice.NewCertificate(repoAll, certPool, log)
	maxmindSvc := expservice.NewMaxmind()
	firewallSvc := expservice.NewFirewall(repoAll, maxmindSvc, log)
	settingSvc := expservice.NewSetting(repoAll, log)
	victoriaMetricsSvc := expservice.NewVictoriaMetrics(repoAll, log)
	fsSvc := expservice.NewFS(repoAll, log)

	shipLog := logger.NewShip(logh)
	brokSH := ship.Default()
	brokSH.Validator = valid
	brokSH.NotFound = shipx.NotFound
	brokSH.HandleError = shipx.HandleError
	brokSH.Logger = shipLog

	{
		healthSvc := brkservice.NewHealth(repoAll, log)
		routes := []shipx.RouteRegister{
			brkrestapi.NewHealth(healthSvc),
		}
		brokRGB := brokSH.Group("/api")
		if err = shipx.RegisterRoutes(brokRGB, routes); err != nil {
			return err
		}
	}

	agentSvc := expservice.NewAgent(repoAll, log)
	agentReleaseSvc := expservice.NewAgentRelease(repoAll, log)
	brokerSvc := expservice.NewBroker(repoAll, log)
	brokerReleaseSvc := expservice.NewBrokerRelease(repoAll, log)
	if err = brokerReset(brokerSvc); err != nil {
		return err
	}

	tunSrvOpts := serverd.Options{
		ConnectListener: tunutil.NewConnectListener(log),
		ConfigLoader:    tunutil.NewAuthConfig(cfg.Database),
		Handler:         brokSH,
		Huber:           hub,
		Validator:       valid.Validate,
		Logger:          log,
		Timeout:         30 * time.Second,
		Context:         ctx,
	}
	brokerTunnelHandler := serverd.NewServer(repoAll, tunSrvOpts)

	jsmodules := []jsvm.Module{
		// jsstd.NewCrontab(crond),
	}
	const apiPath = "/api"
	routes := []shipx.RouteRegister{
		exprestapi.NewAgent(agentSvc),
		exprestapi.NewAgentRelease(agentReleaseSvc, brokerSvc),
		exprestapi.NewBroker(brokerSvc),
		exprestapi.NewBrokerRelease(brokerReleaseSvc, brokerSvc),
		exprestapi.NewCertificate(certificateSvc),
		exprestapi.NewFirewall(firewallSvc),
		exprestapi.NewFS(fsSvc),
		exprestapi.NewLog(logh),
		exprestapi.NewPlay(jsmodules),
		exprestapi.NewReverse(dualDialer),
		exprestapi.NewSetting(settingSvc),
		exprestapi.NewVictoriaMetrics(victoriaMetricsSvc),
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
	outSH.Logger = shipLog

	firewallMiddle := firewalld.New(firewallSvc, log)
	rootRGB := outSH.Group("/").Use(shipx.NewFirewall(firewallMiddle))
	_ = exprestapi.NewStatic(srvCfg.Static).RegisterRoute(rootRGB)
	apiRGB := rootRGB.Group(apiPath).Use(expmiddle.NewWAF(nil))
	if err = shipx.RegisterRoutes(apiRGB, routes); err != nil { // 注册路由
		return err
	}
	log.Info("HTTP 路由注册完毕")

	// 强制要求统一中间件路由信息。
	//if err = expmiddle.CheckRouteInfo(outSH.Routes()); err != nil {
	//	log.Error("路由信息不符合中间件记录格式", "error", err)
	//	return err
	//}

	cronTasks := []cronv3.Tasker{
		crontab.NewMetrics(victoriaMetricsSvc.PushConfig),
	}
	for _, task := range cronTasks {
		if _, err = crond.AddTask(task); err != nil {
			return err
		}
	}

	listenAddr := srvCfg.Addr
	if listenAddr == "" {
		listenAddr = ":443"
	}
	httpTLS := &tls.Config{GetCertificate: certPool.Match}
	quicTLS := &tls.Config{GetCertificate: certPool.Match, NextProtos: []string{"aegis"}, MinVersion: tls.VersionTLS13}
	httpLog := logger.NewV1(slog.New(logger.Skip(logh, 8)))
	srv := &http.Server{
		Addr:      listenAddr,
		Handler:   outSH,
		TLSConfig: httpTLS,
		ErrorLog:  httpLog,
	}
	quicSrv := &quick.QUICx{
		Addr:       listenAddr,
		Accept:     brokerTunnelHandler,
		QUICConfig: &quic.Config{TLSConfig: quicTLS},
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

	return err
}

func disconnectDB(db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = db.Client().Disconnect(ctx)
}

func brokerReset(brk *expservice.Broker) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return brk.Reset(ctx)
}

func listenQUIC(ctx context.Context, errs chan<- error, srv quick.Server) {
	errs <- srv.ListenAndServe(ctx)
}

func listenHTTPS(errs chan<- error, srv *http.Server) {
	errs <- srv.ListenAndServeTLS("", "")
}
