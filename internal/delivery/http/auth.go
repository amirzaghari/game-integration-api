package http

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Username string `json:"username" binding:"required" example:"testuser1"`
	Password string `json:"password" binding:"required" example:"testpass"`
}

func (r *loginRequest) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode((*struct {
		Username string `json:"username"`
		Password string `json:"password"`
	})(r))
}

// Login godoc
// @Summary Authenticate user
// @Tags Auth
// @Description Authenticate user and return JWT token and username
// @Accept json
// @Produce json
// @Param credentials body loginRequest true "User credentials" example({"username": "testuser1", "password": "testpass"})
// @Success 200 {object} map[string]interface{} "{token: string, username: string}"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Router /auth/login [post]
func (h *Handlers) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.AuthUseCase.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "username": req.Username})
}
