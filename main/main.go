package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/xmx/aegis-server/launch"
)

func main() {
	args := os.Args
	set := flag.NewFlagSet(args[0], flag.ExitOnError)
	ver := set.Bool("v", false, "打印版本")
	cfg := set.String("c", "resources/config/application.json", "配置文件")
	_ = set.Parse(args[1:])
	if *ver {
		return
	}

	signals := []os.Signal{syscall.SIGTERM, syscall.SIGHUP, syscall.SIGKILL, syscall.SIGINT}
	ctx, cancel := signal.NotifyContext(context.Background(), signals...)
	defer cancel()

	logOpt := &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}
	log := slog.New(slog.NewJSONHandler(os.Stdout, logOpt))

	if err := launch.Run(ctx, *cfg); err != nil {
		log.Error("服务运行错误", slog.Any("error", err))
	} else {
		log.Info("服务停止运行")
	}
}
