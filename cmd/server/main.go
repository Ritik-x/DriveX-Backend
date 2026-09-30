package main

import (
	"drivex/internal/config"
	"drivex/internal/database"
	"drivex/internal/handlers"
	"drivex/internal/repository"
	"drivex/internal/routes"
	"drivex/internal/services"
	"drivex/internal/storage"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	//load env

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env")
	}
// Load application configuration
	cfg := config.Load()


	// Connect to PostgreSQL
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	defer db.Close()

	log.Println("PostgreSQL connected successfully")
//repositopry

userRepo := repository.NewUserRepository(db)
sessionRepo := repository.NewSessionRepository(db)
folderRepo := repository.NewwFolderRepository(db)
fileRepo := repository.NewFileRepository(db)
fileShareRepo := repository.NewFileShareRepository(db)
//services

authService := services.NewAuthService(userRepo , sessionRepo,cfg.JWTSecret)
folderService := services.NewFolderService(folderRepo)
fileShareServie := services.NewFileSHareService(fileRepo , fileShareRepo,userRepo)
//handler 

authHandler := handlers.NewAuthHandler(authService)
userHandler := handlers.NewUserHandler()

folderHandler := handlers.NewFolderHandler(folderService)
fileShareHandler := handlers.NewFileShareHandler(fileShareServie)



s3Storage , err := storage.NewS3Storage(&cfg)
if err != nil {
	log.Fatal(err)
}
fileService := services.NewFileService(s3Storage,fileRepo )

fileHandler := handlers.NewFileHandler(
	fileService,
)



	router := gin.Default()

		// Routes
	routes.Setup(
	router,
	authHandler,
	userHandler,
	folderHandler,
	fileHandler,
	fileShareHandler,
	cfg.JWTSecret,
)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	log.Println("DriveX server running on :8080")

	if err := router.Run(":8081"); err != nil {
    log.Fatal(err)

}
}