package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"drivex/internal/models"
	"drivex/internal/redis"
	"drivex/internal/repository"
)

type FolderService struct {
	folderRepo *repository.FolderRepo
	redisClient *redis.Client
}

func NewFolderService(
	folderRepo *repository.FolderRepo,
	redisClient *redis.Client,

) *FolderService {
	return &FolderService{
		folderRepo: folderRepo,
		redisClient: redisClient,
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



	parentKey := "root"

	if parentID != nil {
		parentKey = *parentID
	}

	cacheKey := fmt.Sprintf(
		"folders:owner:%s:parent:%s",
		ownerID,
		parentKey,
	)

	var cachedFolders []models.Folder

	err := s.redisClient.GetJSON(
		ctx,
		cacheKey,
		&cachedFolders,
	)

	if err == nil {
		return cachedFolders, nil
	}

	folders, err := s.folderRepo.GetByOwner(
		ctx,
		ownerID,
		parentID,
	)

	if err != nil {
		return nil, err
	}
	_ = s.redisClient.SetJSON(
		ctx,
		cacheKey,
		folders,
		5*time.Minute,
	)

	return folders, nil

	
	
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