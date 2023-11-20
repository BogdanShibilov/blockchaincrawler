package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/bogdanshibilov/blockchaincrawler/internal/auth/auth"
	"github.com/bogdanshibilov/blockchaincrawler/internal/auth/config"
	v1 "github.com/bogdanshibilov/blockchaincrawler/internal/auth/server/v1"
	"github.com/bogdanshibilov/blockchaincrawler/internal/auth/transport"
)

type App struct {
	logger *zap.SugaredLogger
	cfg    *config.Config
}

func New(logger *zap.SugaredLogger, config *config.Config) *App {
	return &App{
		logger: logger,
		cfg:    config,
	}
}

func (a *App) Run() {
	var (
		cfg = a.cfg
		l   = a.logger
	)

	ctx, cancel := context.WithCancel(context.TODO())
	_ = ctx

	userTransport, err := transport.NewUser(cfg.Transport.UserTransport)
	if err != nil {
		l.Panicf("could not create user transport: %v", err)
	}
	authUseCase := auth.New(cfg.Auth, userTransport)

	l.Infof("Starting server on port %v", cfg.Server.Port)
	grpcService := v1.NewService(authUseCase, l)
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
