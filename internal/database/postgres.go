package database

import (
	"context"
	"fmt"

	"github.com/SANEKNAYMCHIK/newsBot/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, cfg *config.Config) (*Postgres, error) {
	strConn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName
	)
	poolCfg, err := pgxpool.ParseConfig(strConn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
}
