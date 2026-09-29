package services

import (
	"context"
	"errors"
	"strings"

	"drivex/internal/models"
	"drivex/internal/repository"
)

type FolderService struct {
	folderRepo *repository.FolderRepo
}

func NewFolderService(
	folderRepo *repository.FolderRepo,
) *FolderService {
	return &FolderService{
		folderRepo: folderRepo,
	}
}


func (s *FolderService) Create(
	ctx context.Context,
	ownerID string,
	parentID *string,
	name string,
) (*models.Folder, error) {

	name = strings.TrimSpace(name)

	if name == "" {
		return nil, errors.New("folder name is required")
	}

	if len(name) > 255 {
		return nil, errors.New("folder name too long")
	}

	return s.folderRepo.Create(
		ctx,
		ownerID,
		parentID,
		name,
	)
}


func (s *FolderService) GetFolders(
	ctx context.Context,
	ownerID string,
	parentID *string,
) ([]models.Folder, error) {

	return s.folderRepo.GetByOwner(
		ctx,
		ownerID,
		parentID,
	)
}


func (s *FolderService) Update(
	ctx context.Context,
	ownerID string,
	id string,
	name string,
) (*models.Folder, error) {

	name = strings.TrimSpace(name)

	if name == "" {
		return nil, errors.New("folder name is required")
	}

	return s.folderRepo.Update(
		ctx,
		id,
		ownerID,
		name,
	)
}





func (s *FolderService) Delete(
	ctx context.Context,
	ownerID string,
	id string,
) error {

	return s.folderRepo.Delete(
		ctx,
		id,
		ownerID,
	)
}