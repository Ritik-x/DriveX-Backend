package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"image/png"
	"log"
	"path/filepath"
	"strings"

	"drivex/internal/queue"
	"drivex/internal/repository"
	"drivex/internal/storage"

	"github.com/disintegration/imaging"
)

type FileWorker struct {
	rabbit   *queue.RabbitMQ
	s3       *storage.S3Storage
	fileRepo *repository.FileRepository
}

func NewFileWorker(
	rabbit *queue.RabbitMQ,
	s3 *storage.S3Storage,
	fileRepo *repository.FileRepository,
) *FileWorker {
	return &FileWorker{
		rabbit:   rabbit,
		s3:       s3,
		fileRepo: fileRepo,
	}
}

func (w *FileWorker) Start() error {

	messages, err := w.rabbit.Ch.Consume(
		"file_processing",
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	for message := range messages {

		var job FileProcessingJob

		err := json.Unmarshal(message.Body, &job)

		if err != nil {
			log.Printf("invalid file p1rocessing job: %v", err)

			// Message is invalid, so don't retry it.
			if err := message.Nack(false, false); err != nil {
				log.Println("failed to NACK message:", err)
			}

			continue
		}

		log.Printf(
			"processing file: %s",
			job.FileID,
		)

		err = w.processImage(
			context.Background(),
			job,
		)

		if err != nil {
			log.Printf(
				"failed to process file %s: %v",
				job.FileID,
				err,
			)

			// Processing failed.
			// false = don't requeue right now.
			if err := message.Nack(false, false); err != nil {
				log.Println("failed to NACK message:", err)
			}

			continue
		}

		// Processing successful.
		if err := message.Ack(false); err != nil {
			log.Println("failed to ACK message:", err)
		}
	}

	return nil
}

func (w *FileWorker) processImage(
	ctx context.Context,
	job FileProcessingJob,
) error {

	// Only process supported image types.
	ext := strings.ToLower(
		filepath.Ext(job.StorageKey),
	)
	mime := strings.ToLower(job.MimeType)

	if ext != ".jpg" &&
		ext != ".jpeg" &&
		ext != ".png" &&
		mime != "image/jpeg" &&
		mime != "image/jpg" &&
		mime != "image/png" {
		log.Printf(
			"skipping thumbnail for unsupported file: %s",
			job.StorageKey,
		)

		return nil
	}

	if ext == "" {
		switch mime {
		case "image/png":
			ext = ".png"
		case "image/jpeg", "image/jpg":
			ext = ".jpg"
		}
	}

	// Download original file from S3.
	object, err := w.s3.GetObject(
		ctx,
		job.StorageKey,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to get original file: %w",
			err,
		)
	}

	defer object.Body.Close()

	// Decode image.
	img, err := imaging.Decode(object.Body)

	if err != nil {
		return fmt.Errorf(
			"failed to decode image: %w",
			err,
		)
	}

	// Resize while maintaining aspect ratio.
	thumbnail := imaging.Fit(
		img,
		300,
		300,
		imaging.Lanczos,
	)

	var buffer bytes.Buffer

	var contentType string

	switch ext {

	case ".jpg", ".jpeg":

		err = jpeg.Encode(
			&buffer,
			thumbnail,
			&jpeg.Options{
				Quality: 80,
			},
		)

		contentType = "image/jpeg"

	case ".png":

		err = png.Encode(
			&buffer,
			thumbnail,
		)

		contentType = "image/png"
	}

	if err != nil {
		return fmt.Errorf(
			"failed to encode thumbnail: %w",
			err,
		)
	}

	// Example:
	// users/123/files/abc.jpg
	//
	// becomes:
	// users/123/thumbnails/abc.jpg
	parts := strings.Split(
		job.StorageKey,
		"/files/",
	)

	if len(parts) != 2 {
		return fmt.Errorf(
			"invalid storage key: %s",
			job.StorageKey,
		)
	}

	thumbnailKey := fmt.Sprintf(
		"%s/thumbnails/%s",
		parts[0],
		filepath.Base(job.StorageKey),
	)

	// Upload thumbnail to S3.
	err = w.s3.PutObject(
		ctx,
		thumbnailKey,
		&buffer,
		contentType,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to upload thumbnail: %w",
			err,
		)
	}

	// Save thumbnail path in PostgreSQL.
	err = w.fileRepo.SetThumbnailKey(
		ctx,
		job.FileID,
		thumbnailKey,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to save thumbnail key: %w",
			err,
		)
	}

	log.Printf(
		"thumbnail created: %s",
		thumbnailKey,
	)

	return nil
}