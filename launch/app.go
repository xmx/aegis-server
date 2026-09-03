package launch

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/xmx/aegis-server/application/config"
	"github.com/xmx/aegis-server/application/echox"
	"github.com/xmx/aegis-server/application/manager/restapi"
	"github.com/xmx/aegis-server/application/manager/service"
	"github.com/xmx/aegis-server/application/nodelink/linkhub"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/aegis-server/library/mongodb"
	"github.com/xmx/aegis-server/library/netutil"
	"github.com/xmx/aegis-server/library/tlscert"
	"github.com/xmx/aegis-server/library/validation"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// runApp 初始化并启动后端服务。
//
//goland:noinspection GoUnhandledErrorResult
func runApp(ctx context.Context, cfg *config.Config) error {
	lmh, tmpc := initLog()
	defer tmpc.Close()
	log := slog.New(lmh)

	valid := validation.New()
	vexts := append(validation.All() /*, validcustom.All()...*/)
	if err := valid.RegisterCustomValidations(vexts); err != nil {
		return err
	}

	dbopts := options.Client().ApplyURI(cfg.Database.URI)
	mdb, err := mongodb.Connect(dbopts)
	if err != nil {
		return err
	}
	db := repository.NewBaseDB(mdb, log)
	if err = db.CreateIndex(ctx); err != nil {
		log.Error("创建索引出错", "err", err)
		return err
	}

	dcfg, err := db.SystemConfig().GetFallback(ctx)
	hs, logc := initLogHandlers(dcfg.Logger)
	defer logc.Close()

	lmh.Replace(hs...)
	_ = tmpc.Close()

	e := echo.New()
	e.Logger = log
	e.Validator = valid

	for k, v := range dcfg.Server.Static {
		if k != "" && v != "" {
			e.Static(k, v)
		}
	}

	proxyURL, _ := url.Parse("socks5://127.0.0.1:41080")
	cli := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	fake := http.NewServeMux()
	mux := linkhub.NewAgentMUX(db, fake, log)
	agentSvc := service.NewAgent(db, log)
	authSvc := service.NewOAuth(db, cli, log)
	otelClientSvc := service.NewOTelClient(db, log)
	userSvc := service.NewUser(db, log)

	const baseURL = "/api"
	apis := []echox.RouteRegister{
		restapi.NewAgent(agentSvc),
		restapi.NewTunnel(mux, log),
		restapi.NewOAuth(authSvc),
		restapi.NewOTelClient(otelClientSvc),
		restapi.NewPTY(),
		restapi.NewSystemBinary(),
		restapi.NewSystemProcess(),
		restapi.NewSystemService(),
		restapi.NewUser(userSvc),
		restapi.NewWebDAV(baseURL, "/"),
	}

	eg := e.Group(baseURL)
	if err = echox.RegisterRoutes(eg, apis); err != nil {
		return err
	}

	selfTLS := tlscert.NewMatch(nil, log) // 临时自签证书
	addr := dcfg.Server.Addr
	if addr == "" {
		addr = "0.0.0.0:8443"
	}
	srv := &http.Server{
		Addr:              addr,
		Handler:           e,
		ReadHeaderTimeout: 10 * time.Second,
		TLSConfig:         &tls.Config{GetCertificate: selfTLS.GetCertificate},
	}

	errch := make(chan error)
	go serveHTTPS(errch, srv, log)
	defer shutdownHTTPS(ctx, srv, 5*time.Second)

	browserURL := &url.URL{Scheme: "https", Host: netutil.ResolvableAddr(addr)}
	log.Info("请使用浏览器访问", "addr", browserURL.String())

	select {
	case err = <-errch:
	case <-ctx.Done():
		err = ctx.Err()
	}

	cause := context.Cause(ctx)
	log.Error("程序停止运行", "err", err, "cause", cause)

	return err
}
