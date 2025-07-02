package http_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	httpdelivery "gameintegrationapi/internal/delivery/http"
	"gameintegrationapi/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockWalletUseCase struct{}

func (m *mockWalletUseCase) Withdraw(userID uint, amount float64, currency, providerTx, roundID, gameID string) (*domain.Transaction, error) {
	return &domain.Transaction{
		ID:           1,
		ProviderTxID: providerTx,
		OldBalance:   1000,
		NewBalance:   900,
		Status:       "PLACED",
	}, nil
}

func (m *mockWalletUseCase) Deposit(userID uint, amount float64, currency, providerTx, providerWithdrawnTxID string) (*domain.Transaction, error) {
	return nil, nil
}
func (m *mockWalletUseCase) Cancel(userID uint, providerTx string) (*domain.Transaction, error) {
	return nil, nil
}

func TestWithdrawRejectsExtraFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &httpdelivery.Handlers{WalletUseCase: &mockWalletUseCase{}}
	r := gin.New()
	r.POST("/bet/withdraw", func(c *gin.Context) {
		c.Set("userID", uint(1))
		h.Withdraw(c)
	})

	body := map[string]interface{}{
		"currency":                "USD",
		"amount":                  100,
		"provider_transaction_id": "provider-tx-123",
		"extra":                   "notallowed",
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/bet/withdraw", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "unknown field")
}

func TestWithdrawAcceptsValidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &httpdelivery.Handlers{WalletUseCase: &mockWalletUseCase{}}
	r := gin.New()
	r.POST("/bet/withdraw", func(c *gin.Context) {
		c.Set("userID", uint(1))
		h.Withdraw(c)
	})

	body := map[string]interface{}{
		"currency":                "USD",
		"amount":                  100,
		"provider_transaction_id": "provider-tx-123",
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/bet/withdraw", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "tx-1")
}
 