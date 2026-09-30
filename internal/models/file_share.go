package models

import "time"

type FileShare struct {
	ID               string    `json:"id"`
	FileID           string    `json:"file_id"`
	OwnerID          string    `json:"owner_id"`
	SharedWithUserID string    `json:"shared_with_user_id"`
	Permission       string    `json:"permission"`
	CreatedAt        time.Time `json:"created_at"`
}