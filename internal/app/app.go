package app

import (
	"context"
	"fmt"

	"github.com/itmo-lite-chat/go-template-svc.git/cmd/config"
	grpcserver "github.com/itmo-lite-chat/go-template-svc.git/internal/app/grpc_server"
	"github.com/itmo-lite-chat/go-utils/closer"
	"github.com/itmo-lite-chat/go-utils/logger"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Server interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type app struct {
	config     *config.Config
	grpcServer Server
	cancel     context.CancelFunc
	closer     *closer.Closer
}

func NewApp(cfg *config.Config) *app {
	return &app{
		config: cfg,
		closer: closer.NewCloser(),
	}
}

func (a *app) Run(c context.Context) error {
	ctx, cancel := context.WithCancel(c)
	a.cancel = cancel

	context.AfterFunc(ctx, func() {
		defer func() {
			if e := recover(); e != nil {
				logger.Error(ctx, "panic: context after-func", "recover", e)
			}
		}()
		if err := a.GracefulShutdown(ctx); err != nil {
			logger.Error(ctx, "graceful shutdown error", err)
		}
	})

	_, err := a.initPostgresDB(ctx)
	if err != nil {
		return errors.Wrap(err, "can't init PostgresDB")
	}

	eg, ctx := errgroup.WithContext(c)

	eg.Go(func() error {
		defer a.cancel()
		return a.startGRPCServer(ctx)
	})

	logger.Info(ctx, "сервер поднялся")

	if err := eg.Wait(); err != nil {
		return errors.Wrap(err, "errorgroup вернула ошибку")
	}
	return nil
}

func (a *app) GracefulShutdown(ctx context.Context) error {
	var err error
	if err = a.grpcServer.Shutdown(ctx); err != nil {
		err = errors.Wrap(err, "can't shutdown grpc server")
	}

	if err = a.closer.Close(); err != nil {
		err = errors.Wrap(err, "can't close all connections")
	}

	return err
}

func (a *app) startGRPCServer(ctx context.Context) error {
	defer func() {
		if e := recover(); e != nil {
			logger.Error(ctx, "panic: grpc start", "panic", e)
		}
	}()

	addr := fmt.Sprintf("%s:%d", a.config.GrpcServer.Host, a.config.GrpcServer.Port)
	s := grpcserver.NewServer(addr)
	logger.Debug(
		ctx, "gRPC server started",
		"addr", addr,
	)

	a.grpcServer = s

	if err := a.grpcServer.Run(ctx); err != nil {
		return fmt.Errorf("grpc server is shutdown: %w", err)
	}
	return nil
}
