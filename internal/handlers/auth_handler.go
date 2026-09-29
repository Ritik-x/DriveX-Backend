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

type LoginRequest struct {
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

func ( h *AuthHandler) Login( c *gin.Context){
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	user , acessToekn , refreshToken , err := h.authService.Login(
		c.Request.Context(), 
		req.Email,
		req.Password,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":          user,
		"access_token":  acessToekn,
		"refresh_token": refreshToken,
	})
}