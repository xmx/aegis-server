package launch

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/robfig/cron/v3"
	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/business/bservice"
	"github.com/xmx/aegis-server/business/service"
	"github.com/xmx/aegis-server/business/validext"
	"github.com/xmx/aegis-server/channel/broker"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/aegis-server/handler/brkapi"
	"github.com/xmx/aegis-server/handler/middle"
	"github.com/xmx/aegis-server/handler/restapi"
	"github.com/xmx/aegis-server/handler/shipx"
	"github.com/xmx/aegis-server/library/cronv3"
	"github.com/xmx/aegis-server/library/validation"
	"github.com/xmx/aegis-server/logger"
	"github.com/xmx/aegis-server/profile"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/connstring"
	"golang.org/x/net/quic"
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
	valid.RegisterCustomValidations(validext.All())
	if err := valid.Validate(cfg); err != nil {
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

	crontab := cronv3.New(cron.WithSeconds())
	crontab.Start()
	defer crontab.Stop()

	mongoDB := cli.Database(mongoURL.Database)
	repoAll := repository.NewAll(mongoDB)

	if err = repoAll.CreateIndex(ctx); err != nil {
		return err
	}
	log.Info("数据库索引建立完毕")

	brokHub := broker.NewHub()
	certificateSvc := service.NewCertificate(repoAll, log)
	termSvc := service.NewTerm(log)

	inSH := ship.Default()
	inSH.Name = "server.internal"
	inSH.Validator = valid
	inSH.NotFound = shipx.NotFound
	inSH.HandleError = shipx.HandleErrorWithHost("server.internal")
	inSH.Logger = logger.NewShip(logHandler, 6)

	{
		aliveSvc := bservice.NewAlive(repoAll, log)
		systemSvc := bservice.NewSystem(repoAll, log)
		routes := []shipx.RouteRegister{
			brkapi.NewAlive(aliveSvc),
			brkapi.NewSystem(systemSvc),
		}
		inRGB := inSH.Group("/api")
		if err = shipx.RegisterRoutes(inRGB, routes); err != nil {
			return err
		}
	}

	brokerSvc := service.NewBroker(repoAll, brokHub, log)
	if err = brokerReset(brokerSvc); err != nil {
		return err
	}

	brokGate := broker.NewGate(repoAll, brokHub, inSH, log)

	const apiPath = "/api"
	routes := []shipx.RouteRegister{
		restapi.NewAuth(),
		restapi.NewBroker(brokerSvc),
		restapi.NewCertificate(certificateSvc),
		restapi.NewLog(logHandler),
		restapi.NewChannel(brokGate),
		restapi.NewDAV(apiPath, "/"),
		restapi.NewSystem(),
		restapi.NewTerm(termSvc),
	}

	srvCfg := cfg.Server
	outSH := ship.Default()
	outSH.Name = "server-expose"
	outSH.Validator = valid
	outSH.NotFound = shipx.NotFound
	outSH.HandleError = shipx.HandleError
	outSH.Logger = logger.NewShip(logHandler, 6)

	rootRGB := outSH.Group("/")
	_ = restapi.NewStatic(srvCfg.Static).RegisterRoute(rootRGB)
	apiRGB := rootRGB.Group(apiPath).Use(middle.WAF(nil))
	if err = shipx.RegisterRoutes(apiRGB, routes); err != nil { // 注册路由
		return err
	}
	log.Info("HTTP 路由注册完毕")

	listenAddr := srvCfg.Addr
	if listenAddr == "" {
		listenAddr = ":https"
	}
	tlsCfg := &tls.Config{
		GetCertificate: certificateSvc.GetCertificate,
		NextProtos:     []string{"h2", "http/1.1", "aegis"},
	}
	httpLog := logger.NewV1(slog.New(logger.Skip(logHandler, 8)))
	srv := &http.Server{
		Addr:      listenAddr,
		Handler:   outSH,
		TLSConfig: tlsCfg,
		ErrorLog:  httpLog,
	}
	errs := make(chan error)
	go listenAndServe(errs, srv, log)
	select {
	case err = <-errs:
	case <-ctx.Done():
	}
	_ = srv.Close()
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

func brokerReset(brk *service.Broker) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return brk.Reset(ctx)
}

func listenAndServe(errs chan<- error, srv *http.Server, log *slog.Logger) {
	lc := new(net.ListenConfig)
	lc.SetMultipathTCP(true)
	ln, err := lc.Listen(context.Background(), "tcp", srv.Addr)
	if err != nil {
		errs <- err
		return
	}
	laddr := ln.Addr().String()
	log.Warn("http 服务监听成功", "listen", laddr)

	errs <- srv.ServeTLS(ln, "", "")
}

func listenQUIC(errs chan<- error, addr string, quicCfg *quic.Config) {
	end, err := quic.Listen("udp", addr, quicCfg)
	if err != nil {
		errs <- err
		return
	}

	end.Accept(context.Background())
}
