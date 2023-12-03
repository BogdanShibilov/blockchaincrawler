package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/metrics"
)

func Metrics() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		path := ctx.Request.URL.Path
		status := strconv.Itoa(ctx.Writer.Status())
		method := ctx.Request.Method
		metrics.HttpResponseTime.WithLabelValues(path, status, method).Observe(time.Since(start).Seconds())
		metrics.HttpRequestsTotalCollector.WithLabelValues(path, status, method).Inc()
	}
}
