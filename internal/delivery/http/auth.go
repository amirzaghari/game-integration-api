package http

import (
	"net/http"
	"os"
	"strings"

	"gameintegrationapi/internal/infrastructure"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginRequest example
// @Description Example: {"username": "testuser1", "password": "testpass"}
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"testuser1"`
	Password string `json:"password" binding:"required" example:"testpass"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// JWT middleware for Gin
func JWTAuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Error: "Missing or invalid token"})
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := infrastructure.ParseJWT(tokenStr, os.Getenv("JWT_SECRET"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid token"})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

// Login authenticates a player attempting to play a game.
// @Summary Authenticate user
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User credentials" example({"username": "testuser1", "password": "testpass"})
// @Success 200 {object} LoginResponse "OK" example({"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."})
// @Failure 400 {object} ErrorResponse "Bad Request" example({"error": "Invalid request"})
// @Failure 401 {object} ErrorResponse "Unauthorized" example({"error": "Invalid credentials"})
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var user struct {
		ID       uint
		Username string
		Password string
	}
	if err := db.Table("users").Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
		return
	}
	token, err := infrastructure.GenerateJWT(user.ID, os.Getenv("JWT_SECRET"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, LoginResponse{Token: token})
}
