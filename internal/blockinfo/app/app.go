package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/blockinfo"
	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/config"
	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/database/postgres"
	v1 "github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/grpcserver/v1"
	"github.com/bogdanshibilov/blockchaincrawler/internal/blockinfo/repository"
	"go.uber.org/zap"
)

type App struct {
	logger *zap.SugaredLogger
	cfg    *config.Config
}

func New(l *zap.SugaredLogger, c *config.Config) *App {
	return &App{
		logger: l,
		cfg:    c,
	}
}

func (a *App) Run() {
	var (
		cfg = a.cfg
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
			l.Panicf("failed to close main db err: %v", err)
		}
		l.Info("main db was succesfully closed")
	}()

	blockRepo := repository.NewBlock(mainDb)
	blockUseCase := blockinfo.New(blockRepo)

	l.Infof("Starting server on port %v", cfg.Server.Port)
	grpcService := v1.NewService(blockUseCase, l)
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
