package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/itmo-lite-chat/go-utils/logger"
	"github.com/itmo-lite-chat/messages_svc/cmd/config"
	"github.com/itmo-lite-chat/messages_svc/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatalf(ctx, "can't parse app config: %v", err)
	}

	level, err := zap.ParseAtomicLevel(cfg.LogLevel)
	if err != nil {
		logger.Fatal(ctx, "can't parse logger level",
			"level", cfg.LogLevel,
			err,
		)
	}
	globalLogger := logger.NewLogger(level)
	logger.SetLogger(globalLogger)

	go signalHandler(ctx, cancel)

	application := app.NewApp(cfg)
	if err = application.Run(ctx); err != nil {
		logger.Error(ctx, "can't run application", err)
		return
	}
	logger.Info(ctx, "application is shutdown normally")
}

func signalHandler(ctx context.Context, cancelFunc context.CancelFunc) {
	osSigCh := make(chan os.Signal, 1)

	signal.Notify(
		osSigCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	defer signal.Stop(osSigCh)
	s := <-osSigCh
	logger.Info(ctx, "получен signal",
		"signal", s,
	)
	cancelFunc()
}
