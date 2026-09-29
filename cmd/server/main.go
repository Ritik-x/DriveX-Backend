package main

import (
	"drivex/internal/config"
	"drivex/internal/database"
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


	router := gin.Default()
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