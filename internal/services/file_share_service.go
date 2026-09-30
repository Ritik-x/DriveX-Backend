package services

import (
	"context"
	"drivex/internal/models"
	"drivex/internal/repository"
	"errors"
)

type FileShareService struct {

	fileRepo *repository.FileRepository
	fileShareRepo *repository.FileShareREpo
	userRepo *repository.UserRepository

}



func NewFileSHareService ( fileRepo *repository.FileRepository,
	fileShareRepo *repository.FileShareREpo,
	userRepo *repository.UserRepository) *FileShareService {

		return &FileShareService{

			fileRepo:      fileRepo,
		fileShareRepo: fileShareRepo,
		userRepo:      userRepo,
		}
	}

	func ( s *FileShareService) ShareFile(

ctx context.Context,
ownerId string , 
fileId string , 
email string , 
permission string ,

	) error {

			if permission != "viewer" && permission != "editor" {
		return errors.New("invalid permission")
	}





// verify for the users
		_, err := s.fileRepo.GetFileById(ctx, fileId, ownerId)

	if err != nil {
		return errors.New("file not found")
	}
	user, err := s.userRepo.GetByEmail(ctx, email)

	if err != nil {
		return errors.New("user not found")
	}

	if user.ID == ownerId {
		return errors.New("cannot share file with owner")
	}
	share := &models.FileShare{
		FileID:           fileId,
		OwnerID:          ownerId,
		SharedWithUserID: user.ID,
		Permission:       permission,
	}
return s.fileShareRepo.Creaate(ctx, share)
	}