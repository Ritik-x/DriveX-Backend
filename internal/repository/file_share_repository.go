package repository

import (
	"context"
	"drivex/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)
type FileShareREpo struct {
db *pgxpool.Pool
}

func NewFileShareRepository(db *pgxpool.Pool) *FileShareREpo {
	return  &FileShareREpo{
		db :db,
	}
}


func ( r *FileShareREpo) Creaate( ctx context.Context , share *models.FileShare) error {
	_ , err := r.db.Exec(ctx , `
		INSERT INTO file_shares (
			file_id,
			owner_id,
			shared_with_user_id,
			permission
		)
		VALUES ($1, $2, $3, $4)
	`, share.FileID ,share.OwnerID,
		share.SharedWithUserID,
		share.Permission,)
		return err
}


func( r *FileShareREpo) GetPermission( ctx context.Context ,  fileId string , ownerId string )([]models.FileShare,error){


	rows , err := r.db.Query( ctx , `  SELECT id , file_id , owner_id , 	shared_with_user_id,
			permission,
			created_at
		FROM file_shares
		WHERE file_id = $1
		  AND owner_id = $2
		ORDER BY created_at DESC ` , fileId,
		ownerId,)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var shares []models.FileShare
	for rows.Next(){
		var share models.FileShare


		err := rows.Scan(
			&share.ID,
			&share.FileID,
			&share.OwnerID,
				&share.SharedWithUserID,
			&share.Permission,
			&share.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		shares = append(shares, share)
	}
	return shares, rows.Err()

}


func (r *FileShareREpo) Delete(
	ctx context.Context,
	fileID string,
	ownerID string,
	sharedWithUserID string,
) error {

	_, err := r.db.Exec(ctx, `
		DELETE FROM file_shares
		WHERE file_id = $1
		  AND owner_id = $2
		  AND shared_with_user_id = $3
	`,
		fileID,
		ownerID,
		sharedWithUserID,
	)

	return err
}



func ( r *FileShareREpo) GetSharedFilewithuser(ctx context.Context , userId string ) ([]models.File , error){
		rows, err := r.db.Query(ctx, `
		SELECT
			f.id,
			f.owner_id,
			f.folder_id,
			f.name,
			f.original_name,
			f.storage_key,
			f.mime_type,
			f.size,
			f.deleted_at,
			f.created_at,
			f.updated_at
		FROM files f
		INNER JOIN file_shares fs
			ON fs.file_id = f.id
		WHERE fs.shared_with_user_id = $1
		  AND f.deleted_at IS NULL
		ORDER BY f.created_at DESC
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()


	var files []models.File

	for rows.Next() {
		var file models.File

		err := rows.Scan(
			&file.FolderId,
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