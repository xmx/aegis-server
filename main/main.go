package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"syscall"

	"github.com/xmx/aegis-server/infra/banner"
	"github.com/xmx/aegis-server/launch"
)

func main() {
	args := os.Args
	name := filepath.Base(os.Args[0])
	set := flag.NewFlagSet(name, flag.ExitOnError)
	ver := set.Bool("v", false, "打印版本")
	cfg := set.String("c", "resources/config/application.jsonc", "配置目录")
	_ = set.Parse(args[1:])
	if banner.ANSI(os.Stdout); *ver {
		return
	}

	if f, _ := os.Create(name + ".crash.txt"); f != nil {
		_ = debug.SetCrashOutput(f, debug.CrashOptions{})
		_ = f.Close()
	}

	signals := []os.Signal{syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT}
	ctx, cancel := signal.NotifyContext(context.Background(), signals...)
	defer cancel()

	opt := &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}
	log := slog.New(slog.NewJSONHandler(os.Stdout, opt))

	if err := launch.Run(ctx, *cfg); err != nil {
		log.Error("服务运行错误", slog.Any("error", err))
	} else {
		log.Info("服务停止运行")
	}
}
