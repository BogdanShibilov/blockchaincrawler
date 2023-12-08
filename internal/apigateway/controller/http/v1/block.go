package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/apigateway"
)

type BlockRoutes struct {
	api apigateway.BlocksUseCase
	l   *zap.SugaredLogger
}

func NewBlockRoutes(handler *gin.RouterGroup,
	api apigateway.BlocksUseCase,
	l *zap.SugaredLogger,
) {

	r := &BlockRoutes{
		api: api,
		l:   l,
	}

	blockHandler := handler.Group("/block")
	{
		blockHandler.GET("/recent", r.GetLastBlocks)
		blockHandler.GET("/header", r.GetHeaders)
		blockHandler.GET("/transaction/:blockhash", r.GetTxsByBlockHash)
		blockHandler.GET("/withdrawal/:blockhash", r.GetWsByBlockHash)
	}
}

// GetHeaders godoc
// @Summary Get headers array
// @Description Returns paginated list of block headers
// @Tags block
// @Produce json
// @Success 200 {object} dto.PagedDto
// @Failure 500
// @Router /block/header [get]
func (r *BlockRoutes) GetHeaders(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	dto, err := r.api.GetHeaders(ctx, page, pageSize)
	if err != nil {
		r.l.Errorf("failed to get headers from blockinfo: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "failed to get headers")
		return
	}

	ctx.JSON(http.StatusOK, dto)
}

// GetTxsByBlockHash godoc
// @Summary Get transactions array
// @Description Returns paginated list of transactions
// @Tags block
// @Produce json
// @Param blockhash path string true "looks up txs by this hash"
// @Success 200 {object} dto.PagedDto
// @Failure 400
// @Failure 500
// @Router /block/transaction/{blockhash} [get]
func (r *BlockRoutes) GetTxsByBlockHash(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	blockHash := ctx.Param("blockhash")
	if blockHash == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &gin.H{"message": "specify hash parameter"})
		return
	}

	dto, err := r.api.GetTxsByBlockHash(ctx, blockHash, page, pageSize)
	if err != nil {
		r.l.Errorf("failed to get transactions from blockinfo: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "failed to get transactions")
		return
	}

	ctx.JSON(http.StatusOK, dto)
}

// GetWsByBlockHash godoc
// @Summary Get withdrawals array
// @Description Returns paginated list of withdrawals
// @Tags block
// @Produce json
// @Param blockhash path string true "looks up withdrawals by this hash"
// @Success 200 {object} dto.PagedDto
// @Failure 400
// @Failure 500
// @Router /block/withdrawal/{blockhash} [get]
func (r *BlockRoutes) GetWsByBlockHash(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	blockHash := ctx.Param("blockhash")
	if blockHash == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &gin.H{"message": "specify hash parameter"})
		return
	}

	dto, err := r.api.GetWsByBlockHash(ctx, blockHash, page, pageSize)
	if err != nil {
		r.l.Errorf("failed to get withdrawals from blockinfo: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "failed to get withdrawals")
		return
	}

	ctx.JSON(http.StatusOK, dto)
}

// GetLastBlocks godoc
// @Summary Get recent blocks
// @Description Returns array of recently discovered blocks
// @Tags block
// @Produce json
// @Success 200 {array} dto.BlockDto
// @Failure 500
// @Router /block/recent [get]
func (r *BlockRoutes) GetLastBlocks(ctx *gin.Context) {
	blocks, err := r.api.GetLastNBlocks(ctx, 20)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "failed to get most recent blocks")
		return
	}

	ctx.JSON(http.StatusOK, blocks)
}
