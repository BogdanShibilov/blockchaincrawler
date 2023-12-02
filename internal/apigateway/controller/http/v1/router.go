package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/apigateway"
	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/config"
)

type Router struct {
	api apigateway.UseCase
	l   *zap.SugaredLogger
	cfg *config.Config
}

func NewRouter(api apigateway.UseCase, l *zap.SugaredLogger, cfg *config.Config) *Router {
	return &Router{
		api: api,
		l:   l,
		cfg: cfg,
	}
}

func (r *Router) Run(handler *gin.Engine) {
	handler.GET("/healthz", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })

	h := handler.Group("/api/v1")
	{
		NewBlockRoutes(h, r.api, r.l)
		NewAuthRoutes(h, r.api, r.l, r.cfg)
		NewUserRoutes(h, r.api, r.l, r.cfg)
		NewAdminRoutes(h, r.api, r.l, r.cfg)
	}
}
