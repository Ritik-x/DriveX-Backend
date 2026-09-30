package services

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"drivex/internal/models"
	"drivex/internal/repository"
	"drivex/internal/storage"

	"github.com/google/uuid"
)

type FileService struct {
	s3       *storage.S3Storage
	fileRepo *repository.FileRepository
}

func NewFileService(
	s3 *storage.S3Storage,
	fileRepo *repository.FileRepository,
) *FileService {
	return &FileService{
		s3:       s3,
		fileRepo: fileRepo,
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

	// Make sure the storage key belongs to this user.
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

	// Verify t
	actualSize, actualMimeType, err :=
		s.s3.GetObjectMetadata(ctx, storageKey)

	if err != nil {
		return nil, errors.New("file does not exist in storage")
	}

	// S3 metadata is more trustworthy than client-provided size.
	if actualMimeType == "" {
		actualMimeType = mimeType
	}

	file := &models.File{
		Id:           fileID,
		OwnerId:      userID,
		FolderId:    *folderID,
		Name:         fileName,
		OriginalName: originalName,
		StorageKey:   storageKey,
		MimeType:     actualMimeType,
		Size:         actualSize,
	}

	// Save only metadata in PostgreSQL.
	return s.fileRepo.Create(ctx, file)
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

	file, err := s.fileRepo.GetFileById(
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