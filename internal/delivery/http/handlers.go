package http

import (
	"gameintegrationapi/internal/infrastructure"
	"gameintegrationapi/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("your-secret-key") // TODO: Move to config

type Handlers struct {
	AuthUseCase   usecase.AuthUseCase
	PlayerUseCase usecase.PlayerUseCase
	WalletUseCase usecase.WalletUseCase
}

func NewHandlers(authUseCase usecase.AuthUseCase, playerUseCase usecase.PlayerUseCase, walletUseCase usecase.WalletUseCase) *Handlers {
	return &Handlers{
		AuthUseCase:   authUseCase,
		PlayerUseCase: playerUseCase,
		WalletUseCase: walletUseCase,
	}
}

// AuthMiddleware checks JWT and sets user context
func (h *Handlers) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}
		tokenString = tokenString[len("Bearer "):]

		claims, err := infrastructure.ParseJWT(tokenString, string(jwtKey))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
