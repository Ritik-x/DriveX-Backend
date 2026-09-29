package repository

import (
	"context"
	"drivex/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FolderRepo struct {
	db *pgxpool.Pool
}

func NewwFolderRepository (db *pgxpool.Pool ) *FolderRepo{
	return &FolderRepo{
		db :db,
	}
}


func (r *FolderRepo) Create(
	ctx context.Context,
	ownerID string,
	parentID *string,
	name string,
) (*models.Folder, error) {

	folder := &models.Folder{}

	err := r.db.QueryRow(ctx, `
		INSERT INTO folders (
			owner_id,
			parent_id,
			name
		)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			owner_id,
			parent_id,
			name,
			created_at,
			updated_at
	`,
		ownerID,
		parentID,
		name,
	).Scan(
		&folder.ID,
		&folder.OwnerID,
		&folder.ParentID,
		&folder.Name,
		&folder.CreatedAt,
		&folder.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return folder, nil
}


func ( r *FolderRepo) GetByOwner( ctx context.Context,
	ownerID string,
	parentID *string,)([]models.Folder , error){
		rows, err := r.db.Query(ctx, `
		SELECT
			id,
			owner_id,
			parent_id,
			name,
			created_at,
			updated_at
		FROM folders
		WHERE owner_id = $1
		  AND parent_id IS NOT DISTINCT FROM $2
		ORDER BY name ASC
	`, ownerID, parentID)

	if err != nil {
		return nil, err
	}


	defer rows.Close()

	var folders []models.Folder

	for rows.Next() {

		var folder models.Folder

		err := rows.Scan(
			&folder.ID,
			&folder.OwnerID,
			&folder.ParentID,
			&folder.Name,
			&folder.CreatedAt,
			&folder.UpdatedAt,
		)


		if err != nil {
			return nil, err
		}
folders = append(folders, folder)






	}

		if err := rows.Err(); err != nil {
		return nil, err
	}

	return folders, nil
}





func ( r *FolderRepo) GetById(

	ctx context.Context,
	id string,
	ownerID string,
)(*models.Folder,error){
	folder := &models.Folder{}

err := r.db.QueryRow(ctx, `
		SELECT
			id,
			owner_id,
			parent_id,
			name,
			created_at,
			updated_at
		FROM folders
		WHERE id = $1
		  AND owner_id = $2
	`,
		id,
		ownerID,
	).Scan(&folder.ID,
		&folder.OwnerID,
		&folder.ParentID,
		&folder.Name,
		&folder.CreatedAt,
		&folder.UpdatedAt,)
		if err != nil {
		return nil, err
	}
	return folder, nil
}



func( r *FolderRepo) Update(
		ctx context.Context,
	id string,
	ownerID string,
	name string,
)(*models.Folder, error){

		folder := &models.Folder{}

		err := r.db.QueryRow(ctx, `
		UPDATE folders
		SET
			name = $1,
			updated_at = NOW()
		WHERE id = $2
		  AND owner_id = $3
		RETURNING
			id,
			owner_id,
			parent_id,
			name,
			created_at,
			updated_at
	`,
		name,
		id,
		ownerID,
	).Scan(&folder.ID,
		&folder.OwnerID,
		&folder.ParentID,
		&folder.Name,
		&folder.CreatedAt,
		&folder.UpdatedAt,)
			if err != nil {
		return nil, err
	}

	return folder, nil
}



func (r *FolderRepo) Delete(
	ctx context.Context,
	id string,
	ownerID string,
) error {

	_, err := r.db.Exec(ctx, `
		DELETE FROM folders
		WHERE id = $1
		  AND owner_id = $2
	`,
		id,
		ownerID,
	)

	return err
}

