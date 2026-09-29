package routes

import (
	"drivex/internal/handlers"

	"github.com/gin-gonic/gin"
)
func Setup (
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
) {
	api := router.Group("/api/v1")

	auth := api.Group("/auth")

	auth.POST("/register", authHandler.Register)

}