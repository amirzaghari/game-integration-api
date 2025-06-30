package http_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	httpdelivery "gameintegrationapi/internal/delivery/http"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockAuthUseCase struct{}

func (m *mockAuthUseCase) Login(username, password string) (string, error) {
	if username == "user" && password == "pass" {
		return "mocktoken", nil
	}
	return "", errors.New("invalid credentials")
}

func TestLoginRejectsExtraFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &httpdelivery.Handlers{AuthUseCase: &mockAuthUseCase{}}
	r := gin.New()
	r.POST("/auth/login", h.Login)

	body := map[string]interface{}{
		"username": "user",
		"password": "pass",
		"extra":    "notallowed",
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "unknown field")
}

func TestLoginAcceptsValidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &httpdelivery.Handlers{AuthUseCase: &mockAuthUseCase{}}
	r := gin.New()
	r.POST("/auth/login", h.Login)

	body := map[string]interface{}{
		"username": "user",
		"password": "pass",
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "mocktoken")
}
