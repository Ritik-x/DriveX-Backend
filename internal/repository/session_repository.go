package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository (db *pgxpool.Pool ) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (r *SessionRepository ) Create(ctx context.Context , userID string,
	refreshTokenHash string,
	expiresAt time.Time,) error{
		_ , err:= r.db.Exec(ctx , `	INSERT INTO sessions (
			user_id,
			refresh_token_hash,
			expires_at
		)
		VALUES ($1, $2, $3)`  , userID,
		refreshTokenHash,
		expiresAt,)
			return err

	}
