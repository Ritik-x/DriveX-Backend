package routes

import (
	"drivex/internal/handlers"
	"drivex/internal/middleware"
	"time"

	redisclient "drivex/internal/redis"

	"github.com/gin-gonic/gin"
)

func Setup(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	folderHandler *handlers.FolderHandler,
	fileHandler *handlers.FileHAndler,
	fileShareHandler *handlers.FileShareHandler,
	redisClient *redisclient.Client,
	jwtSecret string,
) {
	api := router.Group("/api/v1")

	auth := api.Group("/auth")
	authRateLimit := middleware.RateLimit(redisClient, 50, time.Minute)

	auth.POST("/register", authRateLimit, authHandler.Register)
	auth.POST("/login", authRateLimit, authHandler.Login)
	auth.POST("/refresh", authHandler.Refresh)

	auth.POST("/logout", authHandler.Logout)

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret), middleware.RateLimit(redisClient, 200, time.Minute))

	protected.GET("/me", userHandler.Me)
	protected.GET("/search", fileHandler.Search)
	protected.POST("/folders", folderHandler.Create)
	protected.GET("/folders", folderHandler.GetFolders)
	protected.PATCH("/folders/:id", folderHandler.Update)
	protected.DELETE("/folders/:id", folderHandler.Delete)
	protected.GET("/files", fileHandler.GetFiles)
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
	protected.DELETE("/files/:id", fileHandler.DeleteFile)

	protected.GET("/trash", fileHandler.GetTrash)

	protected.POST("/files/:id/restore", fileHandler.RestoreFile)
	protected.DELETE("/files/:id/permanent", fileHandler.PermanentDeleteFile)
	protected.POST(
		"/files/:id/share",
		fileShareHandler.ShareFile,
	)

	protected.GET("/files/shared", fileHandler.GetSharedFie)
}
