package app

import (
	"blockchaincrawler/internal/auth/auth"
	"blockchaincrawler/internal/auth/config"
	v1 "blockchaincrawler/internal/auth/controller/grpc/v1"
	"blockchaincrawler/internal/auth/transport"
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

type App struct {
	logger *zap.SugaredLogger
	config *config.Config
}

func NewApp(logger *zap.SugaredLogger, config *config.Config) *App {
	return &App{
		logger: logger,
		config: config,
	}
}

func (a *App) Run() {
	var (
		cfg = a.config
		l   = a.logger
	)

	l.Infof("starting grpcServer on port%v", cfg.GrpcServer.Port)

	ctx, cancel := context.WithCancel(context.TODO())
	_ = ctx

	userTransport, err := transport.NewUserGrpcTransport(cfg.Transport.UserGrpc)
	if err != nil {
		l.Panicf("failed to create auth grpc transport: %w", err)
	}

	authUsecase := auth.New(cfg.Auth.JwtSecretKey, userTransport)

	grpcService := v1.NewService(authUsecase, l)
	grpcServer := v1.NewServer(cfg.GrpcServer.Port, grpcService)

	err = grpcServer.Start()
	if err != nil {
		l.Panicf("failed to start grpc server error: %w", err)
	}
	defer grpcServer.Close()

	a.gracefulShutdown(cancel)
}

func (a *App) gracefulShutdown(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
	signal.Stop(ch)
	cancel()
}
