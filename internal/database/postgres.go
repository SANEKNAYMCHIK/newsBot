package database

import (
	"context"
	"fmt"
	"log"

	"github.com/SANEKNAYMCHIK/newsBot/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, cfg *config.Config) (*Postgres, error) {
	strConn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	poolCfg, err := pgxpool.ParseConfig(strConn)

	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	postgres := &Postgres{Pool: pool}

	if err := postgres.Migrate(strConn); err != nil {
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
