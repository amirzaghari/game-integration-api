package http

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics godoc
// @Summary Get application metrics
// @Tags Metrics
// @Description Get application metrics in Prometheus format.
// @Produce plain
// @Success 200 {string} string "Prometheus metrics"
// @Router /metrics [get]
func Metrics() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
