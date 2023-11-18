package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/bogdanshibilov/blockchaincrawler/internal/user/config"
	"github.com/bogdanshibilov/blockchaincrawler/internal/user/database/postgres"
	v1 "github.com/bogdanshibilov/blockchaincrawler/internal/user/grpcserver/v1"
	"github.com/bogdanshibilov/blockchaincrawler/internal/user/repository"
	"github.com/bogdanshibilov/blockchaincrawler/internal/user/user"
)

type App struct {
	logger *zap.SugaredLogger
	config *config.Config
}

func New(logger *zap.SugaredLogger, cfg *config.Config) *App {
	return &App{
		logger: logger,
		config: cfg,
	}
}

func (a *App) Run() {
	var (
		cfg = a.config
		l   = a.logger
	)

	ctx, cancel := context.WithCancel(context.TODO())
	_ = ctx

	l.Infof("Connecting to database %v:%v", cfg.Database.MainNode.Host, cfg.Database.MainNode.Port)
	mainDb, err := postgres.NewWithGorm(cfg.Database.MainNode)
	if err != nil {
		l.Panicf("failed connect to main db on '%s:%d': %v", cfg.Database.MainNode.Host, cfg.Database.MainNode.Port, err)
	}
	defer func() {
		if err := mainDb.Close(); err != nil {
			l.Panicf("failed to close main db err: %w", err)
		}
		l.Info("main db was succesfully closed")
	}()

	userRepo := repository.New(mainDb)
	userUseCase := user.New(userRepo)

	l.Infof("Starting server on port %v", cfg.Server.Port)
	grpcService := v1.NewService(userUseCase, l)
	grpcServer := v1.NewServer(":"+cfg.Server.Port, grpcService)

	err = grpcServer.Start()
	if err != nil {
		l.Panicf("failed to start grpc server error: %v", err)
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
