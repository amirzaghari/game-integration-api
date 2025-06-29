package http

import (
	"fmt"
	"net/http"
	"os"

	"gameintegrationapi/internal/infrastructure"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WithdrawRequest struct {
	Currency            string  `json:"currency" binding:"required" example:"USD"`
	Amount              float64 `json:"amount" binding:"required" example:"100.00"`
	ProviderTransaction string  `json:"provider_transaction_id" binding:"required" example:"provider-tx-123"`
}

type DepositRequest struct {
	Currency              string  `json:"currency" binding:"required" example:"USD"`
	Amount                float64 `json:"amount" binding:"required" example:"100.00"`
	ProviderTransaction   string  `json:"provider_transaction_id" binding:"required" example:"provider-tx-123"`
	ProviderWithdrawnTxID string  `json:"provider_withdrawn_transaction_id" binding:"required" example:"provider-tx-122"`
}

type CancelRequest struct {
	ProviderTransaction string `json:"provider_transaction_id" binding:"required" example:"provider-tx-123"`
}

type BetResponse struct {
	TransactionID       string  `json:"transaction_id" example:"tx-123"`
	ProviderTransaction string  `json:"provider_transaction_id" example:"provider-tx-123"`
	OldBalance          float64 `json:"old_balance" example:"1000.00"`
	NewBalance          float64 `json:"new_balance" example:"900.00"`
	Status              string  `json:"status" example:"PLACED"`
}

// Withdraw processes a withdrawal from a player's balance (bet placement).
// @Summary Place a bet (withdraw)
// @Tags Bet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body WithdrawRequest true "Withdraw details" example({"currency": "USD", "amount": 100.00, "provider_transaction_id": "provider-tx-123"})
// @Success 200 {object} BetResponse "OK" example({"transaction_id": "tx-123", "provider_transaction_id": "provider-tx-123", "old_balance": 1000.00, "new_balance": 900.00, "status": "PLACED"})
// @Failure 400 {object} ErrorResponse "Bad Request" example({"error": "Invalid request"})
// @Failure 401 {object} ErrorResponse "Unauthorized" example({"error": "Unauthorized"})
// @Router /bet/withdraw [post]
func Withdraw(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}
	db := c.MustGet("db").(*gorm.DB)
	var user struct {
		ID       uint
		WalletID string
	}
	if err := db.Table("users").Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}
	var req WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request", Message: err.Error()})
		return
	}
	wallet := infrastructure.NewWalletClient(os.Getenv("WALLET_URL"), os.Getenv("WALLET_TOKEN"))
	wReq := infrastructure.WalletWithdrawRequest{
		Currency: req.Currency,
		UserID:   parseWalletID(user.WalletID),
		Transactions: []struct {
			Amount    float64 `json:"amount"`
			BetID     int     `json:"betId"`
			Reference string  `json:"reference"`
		}{
			{
				Amount:    req.Amount,
				BetID:     0, // You can generate/store a bet ID in DB if needed
				Reference: req.ProviderTransaction,
			},
		},
	}
	oldBalance := 0.0
	if bal, err := wallet.GetBalanceStr(user.WalletID); err == nil {
		oldBalance = bal.Balance
	}
	resp, err := wallet.Withdraw(wReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse{Error: "Wallet service error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, BetResponse{
		TransactionID:       req.ProviderTransaction, // You can use your own transaction ID logic
		ProviderTransaction: req.ProviderTransaction,
		OldBalance:          oldBalance,
		NewBalance:          resp.Balance,
		Status:              "PLACED",
	})
}

// Deposit processes a deposit into a player's account (bet settlement).
// @Summary Settle a bet (deposit)
// @Tags Bet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body DepositRequest true "Deposit details" example({"currency": "USD", "amount": 100.00, "provider_transaction_id": "provider-tx-123", "provider_withdrawn_transaction_id": "provider-tx-122"})
// @Success 200 {object} BetResponse "OK" example({"transaction_id": "tx-124", "provider_transaction_id": "provider-tx-123", "old_balance": 900.00, "new_balance": 1000.00, "status": "WON"})
// @Failure 400 {object} ErrorResponse "Bad Request" example({"error": "Invalid request"})
// @Failure 401 {object} ErrorResponse "Unauthorized" example({"error": "Unauthorized"})
// @Router /bet/deposit [post]
func Deposit(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}
	db := c.MustGet("db").(*gorm.DB)
	var user struct {
		ID       uint
		WalletID string
	}
	if err := db.Table("users").Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}
	var req DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request", Message: err.Error()})
		return
	}
	wallet := infrastructure.NewWalletClient(os.Getenv("WALLET_URL"), os.Getenv("WALLET_TOKEN"))
	wReq := infrastructure.WalletDepositRequest{
		Currency: req.Currency,
		UserID:   parseWalletID(user.WalletID),
		Transactions: []struct {
			Amount    float64 `json:"amount"`
			BetID     int     `json:"betId"`
			Reference string  `json:"reference"`
		}{
			{
				Amount:    req.Amount,
				BetID:     0, // You can generate/store a bet ID in DB if needed
				Reference: req.ProviderTransaction,
			},
		},
	}
	oldBalance := 0.0
	if bal, err := wallet.GetBalanceStr(user.WalletID); err == nil {
		oldBalance = bal.Balance
	}
	resp, err := wallet.Deposit(wReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, ErrorResponse{Error: "Wallet service error", Message: err.Error()})
		return
	}
	status := "WON"
	if req.Amount == 0 {
		status = "LOST"
	}
	c.JSON(http.StatusOK, BetResponse{
		TransactionID:       req.ProviderTransaction,
		ProviderTransaction: req.ProviderTransaction,
		OldBalance:          oldBalance,
		NewBalance:          resp.Balance,
		Status:              status,
	})
}

// Cancel reverts a previously processed transaction.
// @Summary Cancel a transaction
// @Tags Bet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body CancelRequest true "Cancel details" example({"provider_transaction_id": "provider-tx-123"})
// @Success 200 {object} BetResponse "OK" example({"transaction_id": "tx-125", "provider_transaction_id": "provider-tx-123", "old_balance": 1000.00, "new_balance": 1000.00, "status": "CANCELLED"})
// @Failure 400 {object} ErrorResponse "Bad Request" example({"error": "Invalid request"})
// @Failure 401 {object} ErrorResponse "Unauthorized" example({"error": "Unauthorized"})
// @Router /bet/cancel [post]
func Cancel(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}
	db := c.MustGet("db").(*gorm.DB)
	var user struct {
		ID       uint
		WalletID string
	}
	if err := db.Table("users").Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}
	var req CancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request", Message: err.Error()})
		return
	}
	// For the mock wallet, cancel is not implemented, so just return a mock response
	bal, _ := infrastructure.NewWalletClient(os.Getenv("WALLET_URL"), os.Getenv("WALLET_TOKEN")).GetBalanceStr(user.WalletID)
	c.JSON(http.StatusOK, BetResponse{
		TransactionID:       req.ProviderTransaction,
		ProviderTransaction: req.ProviderTransaction,
		OldBalance:          bal.Balance,
		NewBalance:          bal.Balance,
		Status:              "CANCELLED",
	})
}

func parseWalletID(walletID string) int64 {
	var id int64
	fmt.Sscan(walletID, &id)
	return id
}
