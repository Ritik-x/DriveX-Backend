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

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
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


func ( h *AuthHandler) Logout( c *gin.Context){
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req) ; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "refresh_token is required",
	})
	return 



}

err := h.authService.Logout(
		c.Request.Context(),
		req.RefreshToken,
)
if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "logout failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logged out successfully",
	})
}

func ( h *AuthHandler) Refresh( c *gin.Context ) {
	var req RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "refresh_token is required",
		})
		return
	}


accessToken, err := h.authService.Refresh(
		c.Request.Context(),
		req.RefreshToken,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}