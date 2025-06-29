package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Healthz godoc
// @Summary Health check
// @Tags Health
// @Success 200 {string} string "ok"
// @Router /healthz [get]
func Healthz(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
