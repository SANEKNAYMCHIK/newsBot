package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, connStr string) (*Postgres, error) {
	poolCfg, err := pgxpool.ParseConfig(connStr)

	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	postgres := &Postgres{Pool: pool}

	if err := postgres.Migrate(connStr); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database connected and migrated successfully")

	return postgres, nil
}

func (p *Postgres) Migrate(strConn string) error {
	migrator, err := NewMigrator(strConn)
	if err != nil {
		return err
	}
	defer migrator.Close()
	return migrator.Up()
}

func (p *Postgres) Close() {
	p.Pool.Close()
}
