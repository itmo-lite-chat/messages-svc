package app

import (
	"context"
	"fmt"

	"github.com/itmo-lite-chat/go-utils/closer"
	"github.com/itmo-lite-chat/go-utils/logger"
	"github.com/itmo-lite-chat/messages_svc/cmd/config"
	"github.com/itmo-lite-chat/messages_svc/internal/api/grpc_api"
	grpcserver "github.com/itmo-lite-chat/messages_svc/internal/app/grpc_server"
	service "github.com/itmo-lite-chat/messages_svc/internal/services/messages_service"
	storage "github.com/itmo-lite-chat/messages_svc/internal/storage/messages_storage"

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
	serviceAPI *grpc_api.API
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

	// 1. Инициализируем только базу (Mongo)
	db, err := a.initMongoDB(ctx)
	if err != nil {
		return errors.Wrap(err, "can't init MongoDB")
	}

	// 2. Создаем Storage
	msgStorage := storage.NewStorage(db)

	// 3. Создаем Service без реалтайм-клиента
	msgService := service.NewService(msgStorage)

	// 4. Передаем сервис в gRPC хендлер
	a.serviceAPI = grpc_api.NewAPI(msgService)

	eg, ctx := errgroup.WithContext(ctx)
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
	if a.grpcServer != nil {
		if err = a.grpcServer.Shutdown(ctx); err != nil {
			err = errors.Wrap(err, "can't shutdown grpc server")
		}
	}

	if closerErr := a.closer.Close(); closerErr != nil {
		err = errors.Wrap(closerErr, "can't close all connections")
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
	a.grpcServer = grpcserver.NewServer(addr, a.serviceAPI)
	logger.Debug(
		ctx, "gRPC server started",
		"addr", addr,
	)

	if err := a.grpcServer.Run(ctx); err != nil {
		return fmt.Errorf("grpc server is shutdown: %w", err)
	}
	return nil
}
