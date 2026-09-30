package models

import "time"

type File struct {
	Id string `json:"id"`
OwnerId string `json:"owner_id"`
FolderId string  `json:"folder_id,omitempty"`
Name         string     `json:"name"`
    OriginalName string     `json:"original_name"`
    StorageKey   string     `json:"storage_key"`
    MimeType     string     `json:"mime_type"`
    Size         int64      `json:"size"`
    DeletedAt    *time.Time `json:"deleted_at,omitempty"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
}

