package launch

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/application/config"
	"github.com/xmx/aegis-server/application/manager/restapi"
	"github.com/xmx/aegis-server/application/manager/service"
	"github.com/xmx/aegis-server/application/muxgate/linkhub"
	"github.com/xmx/aegis-server/application/shipx"
	"github.com/xmx/aegis-server/datalayer/repository"
	"github.com/xmx/aegis-server/library/logger"
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

	sh := ship.Default()
	sh.Validator = valid
	sh.Logger = logger.NewFormat(log.Handler(), 6)
	sh.NotFound = shipx.NotFound
	sh.HandleError = shipx.HandleError
	for k, v := range dcfg.Server.Static {
		if k != "" && v != "" {
			sh.Route(k).Static(v)
		}
	}

	mux := linkhub.NewAgentMUX(log)
	authSvc := service.NewOAuth(db, log)
	userSvc := service.NewUser(db, log)

	apis := []shipx.RouteRegister{
		restapi.NewTunnel(mux, log),
		restapi.NewOAuth(authSvc),
		restapi.NewUser(userSvc),
	}

	base := sh.Group("/api")
	if err = shipx.RegisterRoutes(base, apis); err != nil {
		return err
	}

	selfTLS := tlscert.NewMatch(nil, log) // 临时自签证书
	addr := dcfg.Server.Addr
	if addr == "" {
		addr = "0.0.0.0:8443"
	}
	srv := &http.Server{
		Addr:              addr,
		Handler:           sh,
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
