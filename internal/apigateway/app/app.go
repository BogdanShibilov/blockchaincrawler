package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/apigateway"
	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/config"
	v1 "github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/controller/http/v1"
	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/transport"
	"github.com/bogdanshibilov/blockchaincrawler/pkg/httpserver"
)

type App struct {
	cfg *config.Config
	l   *zap.SugaredLogger
}

func New(l *zap.SugaredLogger, c *config.Config) *App {
	return &App{
		l:   l,
		cfg: c,
	}
}

func (a *App) Run() {
	var (
		cfg = a.cfg
		l   = a.l
	)

	ctx, cancel := context.WithCancel(context.Background())
	_ = ctx

	blockInfoTransport, err := transport.NewBlockInfo(cfg.Transport.BlockInfoTransport)
	if err != nil {
		l.Panicf("failed to create blocks transport: %v", err)
	}

	authTransport, err := transport.NewAuth(&cfg.Transport.AuthTransport)
	if err != nil {
		l.Panicf("failed to create auth transport: %v", err)
	}

	api := apigateway.NewApi(blockInfoTransport, authTransport)

	handler := gin.Default()
	router := v1.NewRouter(api, l)
	router.Run(handler)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.Http.Port))
	defer httpServer.Shutdown()

	a.gracefulShutdown(cancel)
}

func (a *App) gracefulShutdown(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
	signal.Stop(ch)
	cancel()
}
