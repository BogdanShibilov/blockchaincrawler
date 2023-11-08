package app

import (
	"blockchaincrawler/internal/user/config"
	v1 "blockchaincrawler/internal/user/controller/grpc/v1"
	"blockchaincrawler/internal/user/database/postgres"
	"blockchaincrawler/internal/user/repository"
	"blockchaincrawler/internal/user/user"
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

	l.Infof("trying to run grpcServer on port%v", cfg.GrpcServer.Port)

	ctx, cancel := context.WithCancel(context.TODO())
	_ = ctx

	mainDb, err := postgres.NewWithGorm(cfg.Database.Main)
	if err != nil {
		l.Panicf("failed connect to main db on '%s:%d': %w", cfg.Database.Main.Host, cfg.Database.Main.Port, err)
	}
	defer func() {
		if err := mainDb.Close(); err != nil {
			l.Panicf("failed to close main db err: %w", err)
		}
		l.Info("main db was succesfully closed")
	}()

	repo := repository.New(mainDb)

	userUsecase := user.New(repo)

	grpcService := v1.NewService(userUsecase, l)
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
