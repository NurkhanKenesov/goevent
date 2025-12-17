package auth

import (
	"context"
	"goevent/internal/models"
	"goevent/internal/service"
	"goevent/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var authService = service.NewAuthService()

var jwtSecret = "mysecretkey"

func RegisterUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashed,
	}

	id, err := authService.Register(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = id

	c.JSON(http.StatusOK, gin.H{
		"message": "user registered",
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := authService.GetByEmail(context.Background(), req.Email)
	if err != nil || u == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	if !utils.CompareHashAndPassword(u.Password, req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(u.Email, jwtSecret, 72*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
	})
}

func Profile(c *gin.Context) {
	emailVal, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	email, ok := emailVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid context data"})
		return
	}

	u, err := authService.GetByEmail(context.Background(), email)
	if err != nil || u == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       u.ID,
			"username": u.Username,
			"email":    u.Email,
		},
	})
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "server is running"})
}
