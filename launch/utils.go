package launch

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/xmx/aegis-common/logger"
	"github.com/xmx/aegis-server/application/config"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/library/netutil"
	"gopkg.in/natefinch/lumberjack.v2"
)

func serveHTTPS(errch chan error, srv *http.Server, log *slog.Logger) {
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		errch <- err
		log.Error("网络监听错误", "network", "tcp", "addr", addr, "err", err)
		return
	}

	laddr := ln.Addr()
	if a, ok := laddr.(*net.TCPAddr); ok {
		if name := netutil.IsRestrictedPort(uint32(a.Port)); name != "" {
			log.Warn(
				"【温馨提示】当前端口被 Chromium/Chrome 等浏览器标记为受限端口，可能出现 ERR_UNSAFE_PORT，建议更换监听端口",
				"port", a.Port, "service", name,
				"spec", "https://fetch.spec.whatwg.org/#port-blocking",
				"source", "https://chromium.googlesource.com/chromium/src/+/refs/tags/150.0.7865.1/net/base/port_util.cc#30",
			)
		}
	}

	errch <- srv.ServeTLS(ln, "", "")
}

func shutdownHTTPS(parent context.Context, srv *http.Server, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	return srv.Shutdown(ctx)
}

func initLog() (*logger.MultiHandler, io.Closer) {
	opts := &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}
	file := &lumberjack.Logger{Filename: config.DefaultLogFilename}

	return logger.NewMultiHandler(
		logger.NewTint(os.Stdout, opts),
		slog.NewJSONHandler(file, opts),
	), file
}

func initLogHandlers(l model.SCLogger) (hs []slog.Handler, c io.Closer) {
	lvl := l.LevelVar()
	opts := &slog.HandlerOptions{AddSource: true, Level: lvl}
	if fd := l.ConsoleFD(); fd != nil {
		hs = append(hs, logger.NewTint(fd, opts))
	}
	if f := l.Lumber(); f != nil {
		hs = append(hs, slog.NewJSONHandler(f, opts))
		return hs, f
	}

	return hs, io.NopCloser(nil)
}
