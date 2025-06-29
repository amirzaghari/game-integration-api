package http

import (
	"net/http"
	"os"

	"gameintegrationapi/internal/infrastructure"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProfileResponse struct {
	UserID   uint    `json:"user_id" example:"1"`
	Balance  float64 `json:"balance" example:"5000.00"`
	Currency string  `json:"currency" example:"USD"`
}

// Profile retrieves essential player details.
// @Summary Get player profile
// @Tags Player
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ProfileResponse "OK" example({"user_id": 1, "balance": 5000.00, "currency": "USD"})
// @Failure 401 {object} ErrorResponse "Unauthorized" example({"error": "Unauthorized"})
// @Router /profile [get]
func Profile(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}
	db := c.MustGet("db").(*gorm.DB)
	var user struct {
		ID       uint
		WalletID string
		Currency string
	}
	if err := db.Table("users").Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}
	wallet := infrastructure.NewWalletClient(os.Getenv("WALLET_URL"), os.Getenv("WALLET_TOKEN"))
	balanceResp, err := wallet.GetBalanceStr(user.WalletID)
	if err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse{Error: "Wallet service error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ProfileResponse{
		UserID:   user.ID,
		Balance:  balanceResp.Balance,
		Currency: balanceResp.Currency,
	})
}
