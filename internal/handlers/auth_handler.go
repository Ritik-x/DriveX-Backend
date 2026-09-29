package handlers

import (
	"drivex/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler( authService *services.AuthService, ) *AuthHandler{
	return &AuthHandler{
			authService: authService,
	}
}
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ( h *AuthHandler) Register ( c *gin.Context){

	var req RegisterRequest

		if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user, err := h.authService.Register(
		c.Request.Context(),
		req.Name,
		req.Email,
		req.Password,
	)

		if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
c.JSON(http.StatusCreated, gin.H{
		"user": user,
	})
}