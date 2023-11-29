package v1

import (
	"net/http"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/apigateway"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Router struct {
	api apigateway.UseCase
	l   *zap.SugaredLogger
}

func NewRouter(api apigateway.UseCase, l *zap.SugaredLogger) *Router {
	return &Router{
		api: api,
		l:   l,
	}
}

func (r *Router) Run(handler *gin.Engine) {
	handler.GET("/healthz", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })

	h := handler.Group("/api/v1")
	{
		NewBlockRoutes(h, r.api, r.l)
	}
}
