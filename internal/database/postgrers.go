package database

import (
	"context"
	"drivex/internal/config"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	// "golang.org/x/tools/go/cfg"
)

func NewPostgres( cfg config.Config )( *pgxpool.Pool , error ) {
dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
	cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
		return pool, nil

}