package routes

import (
	"drivex/internal/handlers"
	"drivex/internal/middleware"

	"github.com/gin-gonic/gin"
)
func Setup (
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	folderHandler *handlers.FolderHandler,
	fileHandler *handlers.FileHAndler,

	jwtSecret string,
) {
	api := router.Group("/api/v1")

	auth := api.Group("/auth")

	auth.POST("/register", authHandler.Register)
auth.POST("/login", authHandler.Login)
auth.POST("/refresh", authHandler.Refresh)

auth.POST("/logout", authHandler.Logout)
	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret))

	protected.GET("/me", userHandler.Me)
	protected.POST("/folders",folderHandler.Create )
protected.GET("/folders", folderHandler.GetFolders)
protected.PATCH("/folders/:id", folderHandler.Update)
protected.DELETE("/folders/:id", folderHandler.Delete)
protected.POST(
	"/files/upload-url",
	fileHandler.GenerateUpladUrl,
)
// protected.POST(
// 	"/files/upload-url",
// 	fileHandler.GenerateUpladUrl,
// )

protected.POST(
	"/files/complete",
	fileHandler.CompleteUpload,
)
protected.GET(
	"/files/:id/download",
	fileHandler.Download,
)
}