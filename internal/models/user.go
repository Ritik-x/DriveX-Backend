package models

import "time"

type User struct {
	ID string `json:"id"`
	Name string `json:"name"`

Email        string    `json:"email"`
	PasswordHash *string   `json:"-"`
	AvatarURL    *string   `json:"avatar_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
