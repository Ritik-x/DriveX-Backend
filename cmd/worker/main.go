package main

import (
	"log"

	"github.com/joho/godotenv"

	"drivex/internal/config"
	"drivex/internal/database"
	"drivex/internal/queue"
	"drivex/internal/repository"
	"drivex/internal/storage"
	"drivex/internal/workers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found")
	}

	cfg := config.Load()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	log.Println("RabbitMQ URL:", cfg.RabbitMQURL)

	rabbit, err := queue.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbit.Conn.Close()
	defer rabbit.Ch.Close()

	s3Storage, err := storage.NewS3Storage(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := rabbit.DeclareQueue("file_processing"); err != nil {
		log.Fatal(err)
	}

	fileRepo := repository.NewFileRepository(db)
	fileWorker := workers.NewFileWorker(rabbit, s3Storage, fileRepo)

	log.Println("File worker started")

	if err := fileWorker.Start(); err != nil {
		log.Fatal(err)
	}
}