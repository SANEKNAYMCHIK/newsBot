package database

import (
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Migrator struct {
	migrator *migrate.Migrate
}

func NewMigrator(strConn string) (*Migrator, error) {
	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to create source driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, strConn)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return &Migrator{migrator: m}, nil
}

func (m *Migrator) Close() {
	if m.migrator != nil {
		_, _ = m.migrator.Close()
	}
}

func (m *Migrator) Up() error {
	if err := m.migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	log.Println("Migrations completed successfully")
	return nil
}

func (m *Migrator) Down() error {
	if err := m.migrator.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}
	log.Println("All migrations rolled back successfully")
	return nil
}
