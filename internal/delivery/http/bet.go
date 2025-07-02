package http

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type withdrawRequest struct {
	Currency            string  `json:"currency" binding:"required"`
	Amount              float64 `json:"amount" binding:"required"`
	ProviderTransaction string  `json:"provider_transaction_id" binding:"required"`
	RoundID             string  `json:"round_id"`
	GameID              string  `json:"game_id"`
}

func (r *withdrawRequest) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode((*struct {
		Currency            string  `json:"currency"`
		Amount              float64 `json:"amount"`
		ProviderTransaction string  `json:"provider_transaction_id"`
		RoundID             string  `json:"round_id"`
		GameID              string  `json:"game_id"`
	})(r))
}

// Withdraw godoc
// @Summary Place a bet (withdraw)
// @Tags Bet
// @Description Place a bet by withdrawing funds
// @Accept json
// @Produce json
// @Param body body withdrawRequest true "Withdraw details"
// @Success 200 {object} map[string]interface{} "Bet response"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /bet/withdraw [post]
func (h *Handlers) Withdraw(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req withdrawRequest
	if err := c.ShouldBind(&req); err != nil {
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

func (r *depositRequest) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode((*struct {
		Currency              string  `json:"currency"`
		Amount                float64 `json:"amount"`
		ProviderTransaction   string  `json:"provider_transaction_id"`
		ProviderWithdrawnTxID string  `json:"provider_withdrawn_transaction_id"`
	})(r))
}

// Deposit godoc
// @Summary Settle a bet (deposit)
// @Tags Bet
// @Description Settle a bet by depositing funds
// @Accept json
// @Produce json
// @Param body body depositRequest true "Deposit details"
// @Success 200 {object} map[string]interface{} "Bet response"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /bet/deposit [post]
func (h *Handlers) Deposit(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req depositRequest
	if err := c.ShouldBind(&req); err != nil {
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

func (r *cancelRequest) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode((*struct {
		ProviderTransaction string `json:"provider_transaction_id"`
	})(r))
}

// Cancel godoc
// @Summary Cancel a transaction
// @Tags Bet
// @Description Cancel a bet transaction
// @Accept json
// @Produce json
// @Param body body cancelRequest true "Cancel details"
// @Success 200 {object} map[string]interface{} "Bet response"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /bet/cancel [post]
func (h *Handlers) Cancel(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req cancelRequest
	if err := c.ShouldBind(&req); err != nil {
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
