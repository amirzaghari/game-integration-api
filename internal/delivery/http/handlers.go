package http

import (
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	// Add dependencies here (e.g., usecases, logger)
}

// AuthMiddleware checks JWT and sets user context
func (h *Handlers) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement JWT validation and set user context
		c.Next()
	}
}

func (h *Handlers) Healthz(c *gin.Context) {
	c.String(200, "ok")
}

func (h *Handlers) Login(c *gin.Context)    {}
func (h *Handlers) Profile(c *gin.Context)  {}
func (h *Handlers) Withdraw(c *gin.Context) {}
func (h *Handlers) Deposit(c *gin.Context)  {}
func (h *Handlers) Cancel(c *gin.Context)   {}
