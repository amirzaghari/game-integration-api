package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(handlers *Handlers) *gin.Engine {
	r := gin.Default()

	// Redirect root to Swagger UI
	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/auth/login", handlers.Login)
	r.GET("/profile", handlers.AuthMiddleware(), handlers.Profile)
	r.POST("/bet/withdraw", handlers.AuthMiddleware(), handlers.Withdraw)
	r.POST("/bet/deposit", handlers.AuthMiddleware(), handlers.Deposit)
	r.POST("/bet/cancel", handlers.AuthMiddleware(), handlers.Cancel)

	r.GET("/healthz", handlers.Healthz)

	return r
}
