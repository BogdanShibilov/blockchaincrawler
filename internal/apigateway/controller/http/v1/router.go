package v1

import (
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "github.com/bogdanshibilov/blockchaincrawler/docs"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/apigateway"
	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/config"
	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/controller/http/middleware"
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

// @title Blockchain crawler api
// @version 1.0
// @description Allows to get crawled information about blockchain blocks
// @contact.name Bogdan Shibilov
// @contact.email bogdanshibilov@gmail.com
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey	BearerAuth
// @type apiKey
// @name Authorization
// @in header
func (r *Router) Run(handler *gin.Engine) {
	handler.GET("/healthz", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })

	pprof.Register(handler)
	handler.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	h := handler.Group("/api/v1")
	{
		handler.Use(middleware.Metrics())
		NewBlockRoutes(h, r.api, r.l)
		NewAuthRoutes(h, r.api, r.l, r.cfg)
		NewUserRoutes(h, r.api, r.l, r.cfg)
		NewAdminRoutes(h, r.api, r.l, r.cfg)
	}

	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
