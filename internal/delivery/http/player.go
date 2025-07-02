package http

import (
	"net/http"
	"strings"

	"gameintegrationapi/internal/infrastructure"

	"github.com/gin-gonic/gin"
)

// Profile godoc
// @Summary Get player profile
// @Tags Player
// @Description Get the authenticated player's profile
// @Produce json
// @Success 200 {object} map[string]interface{} "Profile response"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
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
		"user_id":  user.ID,
		"balance":  user.Balance,
		"currency": user.Currency,
	})
}
