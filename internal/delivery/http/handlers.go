package http

import (
	"gameintegrationapi/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		id, err := strconv.ParseUint(claims.Subject, 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id in token"})
			c.Abort()
			return
		}

		c.Set("userID", uint(id))
		c.Next()
	}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handlers) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.AuthUseCase.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handlers) Profile(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.PlayerUseCase.GetPlayerInfo(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get user info"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user_id":  user.ID,
		"balance":  user.Balance,
		"currency": user.Currency,
	})
}

type withdrawRequest struct {
	Currency            string  `json:"currency" binding:"required"`
	Amount              float64 `json:"amount" binding:"required"`
	ProviderTransaction string  `json:"provider_transaction_id" binding:"required"`
	RoundID             string  `json:"round_id"`
	GameID              string  `json:"game_id"`
}

func (h *Handlers) Withdraw(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req withdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.WalletUseCase.Withdraw(userID.(uint), req.Amount, req.Currency, req.ProviderTransaction, req.RoundID, req.GameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction_id":          tx.ID,
		"provider_transaction_id": tx.ProviderTxID,
		"old_balance":             tx.OldBalance,
		"new_balance":             tx.NewBalance,
		"status":                  tx.Status,
	})
}

type depositRequest struct {
	Currency              string  `json:"currency" binding:"required"`
	Amount                float64 `json:"amount"`
	ProviderTransaction   string  `json:"provider_transaction_id" binding:"required"`
	ProviderWithdrawnTxID string  `json:"provider_withdrawn_transaction_id" binding:"required"`
}

func (h *Handlers) Deposit(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req depositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.WalletUseCase.Deposit(userID.(uint), req.Amount, req.Currency, req.ProviderTransaction, req.ProviderWithdrawnTxID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction_id":          tx.ID,
		"provider_transaction_id": tx.ProviderTxID,
		"old_balance":             tx.OldBalance,
		"new_balance":             tx.NewBalance,
		"status":                  tx.Status,
	})
}

type cancelRequest struct {
	ProviderTransaction string `json:"provider_transaction_id" binding:"required"`
}

func (h *Handlers) Cancel(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req cancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.WalletUseCase.Cancel(userID.(uint), req.ProviderTransaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction_id":          tx.ID,
		"provider_transaction_id": tx.ProviderTxID,
		"old_balance":             tx.OldBalance,
		"new_balance":             tx.NewBalance,
		"status":                  tx.Status,
	})
}
