package v1

import (
	"net/http"
	"strconv"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/apigateway"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BlockRoutes struct {
	api apigateway.UseCase
	l   *zap.SugaredLogger
}

func NewBlockRoutes(handler *gin.RouterGroup, api apigateway.UseCase, l *zap.SugaredLogger) {
	r := &BlockRoutes{
		api: api,
		l:   l,
	}

	blockHandler := handler.Group("/block")
	{
		blockHandler.GET("/header", r.GetHeaders)
	}
}

func (r *BlockRoutes) GetHeaders(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize"))

	dto, err := r.api.GetHeaders(ctx, page, pageSize)
	if err != nil {
		r.l.Errorf("failed to get headers from blockinfo: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "failed to get headers")
		return
	}

	ctx.JSON(http.StatusOK, dto)
}
