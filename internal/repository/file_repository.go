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

func (r *FileRepository) GetByOwner(
	ctx context.Context,
	ownerID string,
	folderID *string,
) ([]models.File, error) {

	rows, err := r.db.Query(ctx, `
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
		WHERE owner_id = $1
		  AND folder_id IS NOT DISTINCT FROM $2
		  AND deleted_at IS NULL
		ORDER BY created_at DESC
	`,
		ownerID,
		folderID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var files []models.File

	for rows.Next() {

		var file models.File

		err := rows.Scan(
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

		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func ( r *FileRepository) SoftDelete(ctx context.Context , fileId string , ownerId string) error{
		_, err := r.db.Exec(ctx, `
		UPDATE files
		SET deleted_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
		  AND owner_id = $2
		  AND deleted_at IS NULL
	`, fileId, ownerId)

	return err
}

func (r *FileRepository) GetTrash(
	ctx context.Context,
	ownerID string,
) ([]models.File, error) {

	rows, err := r.db.Query(ctx, `
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
		WHERE owner_id = $1
		  AND deleted_at IS NOT NULL
		ORDER BY deleted_at DESC
	`, ownerID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.File

	for rows.Next() {
		var file models.File

		err := rows.Scan(
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

		files = append(files, file)
	}

	return files, rows.Err()
}

func (r *FileRepository) Restore(
	ctx context.Context,
	fileID string,
	ownerID string,
) error {
	_, err := r.db.Exec(ctx, `
		UPDATE files
		SET deleted_at = NULL,
		    updated_at = NOW()
		WHERE id = $1
		  AND owner_id = $2
		  AND deleted_at IS NOT NULL
	`, fileID, ownerID)

	return err
}