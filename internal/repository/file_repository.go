package repository

import (
	"context"
	"drivex/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FileRepository struct {
	db *pgxpool.Pool
}

func NewFileRepository(db *pgxpool.Pool) *FileRepository {
	return &FileRepository{
		db: db,
	}
}

func (r *FileRepository) Create(
	ctx context.Context,
	file *models.File,
) (*models.File, error) {

	result := &models.File{}

	err := r.db.QueryRow(ctx, `
		INSERT INTO files (
		id,
			owner_id,
			folder_id,
			name,
			original_name,
			storage_key,
			mime_type,
			size
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7 ,$8)
		RETURNING
			id,
			owner_id,
			folder_id,
			name,
			original_name,
			storage_key,
			mime_type,
			size,
			deleted_at,
			created_at,
			updated_at
	`,file.Id,
		file.OwnerId,
		file.FolderId,
		file.Name,
		file.OriginalName,
		file.StorageKey,
		file.MimeType,
		file.Size,
	).Scan(
		&result.Id,
		&result.OwnerId,
		&result.FolderId,
		&result.Name,
		&result.OriginalName,
		&result.StorageKey,
		&result.MimeType,
		&result.Size,
		&result.DeletedAt,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}


func ( r *FileRepository) GetFileById( ctx context.Context , id string , ownerId string )(*models.File, error){
		file := &models.File{}

		err := r.db.QueryRow(ctx, `
		SELECT
			id,
			owner_id,
			folder_id,
			name,
			original_name,
			storage_key,
			mime_type,
			size,
			deleted_at,
			created_at,
			updated_at
		FROM files
		WHERE id = $1
		  AND owner_id = $2
		  AND deleted_at IS NULL
	`,
		id,
		file.OwnerId,
	).Scan(
		&file.Id,
		&file.OwnerId,
		&file.FolderId,
		&file.Name,
		&file.OriginalName,
		&file.StorageKey,
		&file.MimeType,
		&file.Size,
		&file.DeletedAt,
		&file.CreatedAt,
		&file.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return file, nil
}