package launch

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/xmx/aegis-server/application/config"
	"github.com/xmx/aegis-server/application/oobe/middle"
	"github.com/xmx/aegis-server/application/oobe/restapi"
	"github.com/xmx/aegis-server/application/oobe/webui"
	"github.com/xmx/aegis-server/library/netutil"
	"github.com/xmx/aegis-server/library/tlscert"
	"github.com/xmx/aegis-server/library/validation"
)

// runOOBE 运行开箱体验，可以通过浏览器进行初始化配置。
//
//goland:noinspection GoUnhandledErrorResult
func runOOBE(ctx context.Context, file string) (*config.Config, error) {
	lmh, cls := initLog()
	defer cls.Close()

	log := slog.New(lmh)
	tmp := make([]byte, 10)
	if _, err := rand.Read(tmp); err != nil {
		log.Error("生成访问凭证错误", "err", err)
		return nil, err
	}
	token := hex.EncodeToString(tmp)

	valid := validation.New()
	vexts := append(validation.All() /*, validcustom.All()...*/)
	if err := valid.RegisterCustomValidations(vexts); err != nil {
		log.Error("校验器注册错误", "err", err)
		return nil, err
	}

	addr, dist := os.Getenv(config.EnvironOOBEAddr), os.Getenv(config.EnvironOOBEDist)
	if addr == "" {
		addr = "0.0.0.0:8443"
		log.Info("如需自定义初始化地址，请通过环境变量设置", "key", config.EnvironOOBEAddr)
	}

	e := echo.New()
	e.Validator = valid
	e.Logger = log

	if dist == "" {
		e.StaticFS("/", webui.Dist())
		log.Info("如需自定义初始化页面，请通过环境变量设置", "key", config.EnvironOOBEDist)
	} else {
		e.Static("/", dist)
	}

	cfgch := make(chan *config.Config)
	oobeAPI := restapi.NewOOBE(cfgch, file)

	eg := e.Group("/api", middle.NewAuth(token))
	_ = oobeAPI.RegisterRoute(eg)

	selfTLS := tlscert.NewMatch(nil, log) // 临时自签证书
	srv := &http.Server{
		Addr:         addr,
		Handler:      e,
		TLSConfig:    &tls.Config{GetCertificate: selfTLS.GetCertificate},
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	errch := make(chan error, 1)
	go serveHTTPS(errch, srv, log)
	defer shutdownHTTPS(ctx, srv, 3*time.Second)

	readableAddr := netutil.ResolvableAddr(addr)
	browserURL := &url.URL{Scheme: "https", Host: readableAddr}
	log.Info("请使用浏览器初始化", "addr", browserURL.String(), "token", token)

	select {
	case cfg := <-cfgch:
		log.Info("初始化完毕")
		return cfg, nil
	case err := <-errch:
		log.Error("初始化服务错误", "err", err)
		return nil, err
	case <-ctx.Done():
		err := ctx.Err()
		cause := context.Cause(ctx)
		log.Error("服务停止运行", "err", err, "cause", cause)
		return nil, err
	}
}
