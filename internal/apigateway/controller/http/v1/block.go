package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/apigateway"
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
		blockHandler.GET("/transaction", r.GetTxsByBlockHash)
		blockHandler.GET("/withdrawal", r.GetWsByBlockHash)
	}
}

func (r *BlockRoutes) GetHeaders(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize"))

	dto, err := r.api.GetHeaders(ctx, page, pageSize)
	if err != nil {
		r.l.Errorf("failed to get headers from blockinfo: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "failed to get headers")
		return
	}

	ctx.JSON(http.StatusOK, dto)
}

func (r *BlockRoutes) GetTxsByBlockHash(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize"))
	blockHash := ctx.Query("hash")

	dto, err := r.api.GetTxsByBlockHash(ctx, blockHash, page, pageSize)
	if err != nil {
		r.l.Errorf("failed to get transactions from blockinfo: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "failed to get transactions")
		return
	}

	ctx.JSON(http.StatusOK, dto)
}

func (r *BlockRoutes) GetWsByBlockHash(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize"))
	blockHash := ctx.Query("hash")

	dto, err := r.api.GetWsByBlockHash(ctx, blockHash, page, pageSize)
	if err != nil {
		r.l.Errorf("failed to get withdrawals from blockinfo: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "failed to get withdrawals")
		return
	}

	ctx.JSON(http.StatusOK, dto)
}
