package main

import (
	"context"
	"drivex/internal/config"
	"drivex/internal/database"
	"drivex/internal/handlers"
	"drivex/internal/repository"
	"drivex/internal/routes"
	"drivex/internal/services"
	"drivex/internal/storage"
	"log"

	redisclient "drivex/internal/redis"

	"drivex/internal/queue"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	//load env
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found, using environment variables")
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

	redisClient := redisclient.NewRedisClient(cfg.RedisAddr)

	if err := redisClient.Ping(context.Background()); err != nil {
		log.Fatal("failed to connect to Redis:", err)
	}

	log.Println("Redis connected")

	rabbit, err := queue.NewRabbitMQ(cfg.RabbitMQURL)

	if err != nil {
		log.Fatal(err)
	}

	defer rabbit.Conn.Close()
	defer rabbit.Ch.Close()

	log.Println("RabbitMQ connected")
	//repositopry

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	folderRepo := repository.NewwFolderRepository(db)
	fileRepo := repository.NewFileRepository(db)
	fileShareRepo := repository.NewFileShareRepository(db)
	//services

	authService := services.NewAuthService(userRepo, sessionRepo, cfg.JWTSecret)
	folderService := services.NewFolderService(folderRepo, redisClient)
	fileShareServie := services.NewFileSHareService(fileRepo, fileShareRepo, userRepo)
	//handler

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler()

	folderHandler := handlers.NewFolderHandler(folderService)
	fileShareHandler := handlers.NewFileShareHandler(fileShareServie)

	s3Storage, err := storage.NewS3Storage(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	fileService := services.NewFileService(s3Storage, fileRepo, rabbit)

	fileHandler := handlers.NewFileHandler(
		fileService,
		folderService,
	)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://drivex-nu.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	// Routes
	routes.Setup(
		router,
		authHandler,
		userHandler,
		folderHandler,
		fileHandler,
		fileShareHandler,
		redisClient,
		cfg.JWTSecret,
	)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	log.Println("DriveX server running on :8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)

	}
}
