package services

import (
	"context"
	"drivex/internal/storage"
	"fmt"
	"path/filepath"
	"uuid"
)

type FileService struct {
	s3 *storage.S3Storage
}
func NewFileService(s3 *storage.S3Storage,) *FileService {
	return &FileService{
		s3 :s3,
	}
}


func (s *FileService) GenerateUploadURL(

ctx context.Context,
userId string ,
fileName string ,
contentType string ,

)(string, string, error) {
	fileId := uuid.New().String()


	extension := filepath.Ext(fileName)

		storageKey := fmt.Sprintf(
		"users/%s/files/%s%s",
		userId,
		fileId,
		extension,
	)
url, err := s.s3.GenerateUploadUrl(
		ctx,
		storageKey,
		contentType,
	)
if err != nil {
		return "", "", err
	}

	return url, storageKey, nil
}