package http

import (
	"net/http"
	"strings"

	"gameintegrationapi/internal/infrastructure"

	"github.com/gin-gonic/gin"
)

type ProfileResponse struct {
	UserID   uint    `json:"user_id" example:"1"`
	Balance  float64 `json:"balance" example:"100.0"`
	Currency string  `json:"currency" example:"USD"`
}

type ProfileErrorResponse struct {
	Error string `json:"error" example:"unauthorized"`
}

// Profile godoc
// @Summary Get player profile
// @Tags Player
// @Description Get the authenticated player's profile
// @Produce json
// @Success 200 {object} ProfileResponse "Profile response"
// @Failure 401 {object} ProfileErrorResponse "Unauthorized"
// @Security BearerAuth
// @Router /profile [get]
func (h *Handlers) Profile(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.PlayerUseCase.GetPlayerInfo(userID.(uint))
	if err != nil {
		if err == infrastructure.ErrWalletUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if strings.Contains(err.Error(), infrastructure.ErrWalletServiceBadRequest.Error()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get user info"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user_id":  user.WalletID,
		"balance":  user.Balance,
		"currency": user.Currency,
	})
}
