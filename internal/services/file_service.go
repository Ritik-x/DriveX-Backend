package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"drivex/internal/models"
	"drivex/internal/queue"
	"drivex/internal/repository"
	"drivex/internal/storage"

	"github.com/google/uuid"
)

type FileService struct {
	s3       *storage.S3Storage
	fileRepo *repository.FileRepository
	rabbit *queue.RabbitMQ
}
type FileProcessingJob struct {
	FileID     string `json:"file_id"`
	StorageKey string `json:"storage_key"`
	MimeType   string `json:"mime_type"`
}
func NewFileService(
	s3 *storage.S3Storage,
	fileRepo *repository.FileRepository,
	rabbit *queue.RabbitMQ,
) *FileService {
	return &FileService{
		s3:       s3,
		fileRepo: fileRepo,
		rabbit:   rabbit,
	}
}

// GenerateUploadURL generates a pre-signed S3 upload URL.
func (s *FileService) GenerateUploadURL(
	ctx context.Context,
	userID string,
	fileName string,
	contentType string,
) (string, string, string, error) {

	// Generate a unique ID for the file.
	fileID := uuid.New().String()

	// Get file extension, e.g. ".pdf", ".jpg"
	extension := filepath.Ext(fileName)

	// Create the S3 object key.
	storageKey := fmt.Sprintf(
		"users/%s/files/%s%s",
		userID,
		fileID,
		extension,
	)

	// Generate pre-signed URL.
	url, err := s.s3.GenerateUploadUrl(
		ctx,
		storageKey,
		contentType,
	)
	if err != nil {
		return "", "", "", err
	}

	return url, storageKey, fileID, nil
}

// func (s *FileService) CompleteUpload(
// 	ctx context.Context,
// 	userID string,
// 	fileID string,
// 	folderID *string,
// 	fileName string,
// 	originalName string,
// 	storageKey string,
// 	mimeType string,
// 	size int64,
//    ) (*models.File, error) {



// 	// Make sure the storage key belongs to this user.
// 	if !isValidStorageKey(userID, storageKey) {
// 		return nil, errors.New("invalid storage key")
// 	}

	
// 	expectedPrefix := fmt.Sprintf(
// 		"users/%s/files/%s",
// 		userID,
// 		fileID,
// 	)

// 	if !strings.HasPrefix(storageKey, expectedPrefix) {
// 		return nil, errors.New("invalid file storage key")
// 	}

// 	// Verify t
// 	actualSize, actualMimeType, err :=
// 		s.s3.GetObjectMetadata(ctx, storageKey)

// 	if err != nil {
// 		return nil, errors.New("file does not exist in storage")
// 	}

// 	// S3 metadata is more trustworthy than client-provided size.
// 	if actualMimeType == "" {
// 		actualMimeType = mimeType
// 	}

// file := &models.File{
// 	Id:           fileID,
// 	OwnerId:      userID,
// 	FolderId:     *folderID,
// 	Name:         fileName,
// 	OriginalName: originalName,
// 	StorageKey:   storageKey,
// 	MimeType:     actualMimeType,
// 	Size:         actualSize,
// }

// err = s.fileRepo.Create(ctx, file)
// if err != nil {
// 	return nil, err
// }

// // Create background processing job.
// job := FileProcessingJob{
// 	FileID:     file.Id,
// 	StorageKey: file.StorageKey,
// 	MimeType:   file.MimeType,
// }

// body, err := json.Marshal(job)
// if err != nil {
// 	return nil, err
// }

// err = s.rabbit.Publish(
// 	"file_processing",
// 	body,
// )

// if err != nil {
// 	return nil, err
// }

// return file, nil
// }







func (s *FileService) CompleteUpload(
	ctx context.Context,
	userID string,
	fileID string,
	folderID *string,
	fileName string,
	originalName string,
	storageKey string,
	mimeType string,
	size int64,
) (*models.File, error) {

	if !isValidStorageKey(userID, storageKey) {
		return nil, errors.New("invalid storage key")
	}

	expectedPrefix := fmt.Sprintf(
		"users/%s/files/%s",
		userID,
		fileID,
	)

	if !strings.HasPrefix(storageKey, expectedPrefix) {
		return nil, errors.New("invalid file storage key")
	}

	actualSize, actualMimeType, err :=
		s.s3.GetObjectMetadata(ctx, storageKey)

	if err != nil {
		return nil, errors.New("file does not exist in storage")
	}

	if actualMimeType == "" {
		actualMimeType = mimeType
	}

	file := &models.File{
		Id:           fileID,
		OwnerId:      userID,
		FolderId:     *folderID,
		Name:         fileName,
		OriginalName: originalName,
		StorageKey:   storageKey,
		MimeType:     actualMimeType,
		Size:         actualSize,
	}

	// 1. Save metadata in PostgreSQL
createdFile, err := s.fileRepo.Create(ctx, file)
if err != nil {
	return nil, err
}

	// 2. Create background job
	job := FileProcessingJob{
		FileID:     file.Id,
		StorageKey: file.StorageKey,
		MimeType:   file.MimeType,
	}

	// 3. Convert job to JSON
	body, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}

	// 4. Publish job to RabbitMQ
	err = s.rabbit.Publish(
		"file_processing",
		body,
	)

	if err != nil {
		return nil, err
	}

	return createdFile, nil
}













// Checks whether the S3 key belongs to the authenticated user.
func isValidStorageKey(
	userID string,
	storageKey string,
) bool {

	prefix := fmt.Sprintf(
		"users/%s/files/",
		userID,
	)

	return strings.HasPrefix(storageKey, prefix)
}



func (s *FileService) GenerateDownloadUrl(ctx context.Context , userId string , fileId string ) (string , error){

	file, err := s.fileRepo.GetIdForUser(
		ctx,
		fileId,
		userId,
	)
	if err != nil {
		return "", errors.New("file not found")
	}

	return s.s3.GenerateDownloadUrl(
		ctx,
		file.StorageKey,
	)



}


func ( s *FileService) GetFiles( ctx context.Context , userId string , folderId *string )([]models.File, error) {

	return s.fileRepo.GetByOwner(
		ctx,
		userId,
		folderId,
	)
}

func( s *FileService) DeletFile( ctx context.Context , fileId string,userId string , ) error{

	return s.fileRepo.SoftDelete(ctx , fileId , userId)
}


func (s *FileService) GetTrash(
	ctx context.Context,
	userID string,
) ([]models.File, error) {
	return s.fileRepo.GetTrash(ctx, userID)
}

func (s *FileService) RestoreFile(
	ctx context.Context,
	fileID string,
	userID string,
) error {
	return s.fileRepo.Restore(ctx, fileID, userID)
}

func (s *FileService) GetSharedFiles(
	ctx context.Context,
	userID string,
) ([]models.File, error) {

	return s.fileRepo.GetSharedFilewithuser(ctx, userID)
}

