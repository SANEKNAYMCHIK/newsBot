package database

import (
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Migrator struct {
	migrator *migrate.Migrate
}

func NewMigrator(strConn string) (*Migrator, error) {
	m, err := migrate.New("iofs://migrations", strConn)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}
	return &Migrator{migrator: m}, nil
}

func (m *Migrator) Close() {
	if m.migrator != nil {
		m.Close()
	}
}

func (m *Migrator) Up() error {
	if err := m.migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	version, dirty, err := m.migrator.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if dirty {
		log.Printf("Database is dirty at version %d", version)
	} else if err == migrate.ErrNilVersion {
		log.Printf("Database is empty, initial state")
	} else {
		log.Printf("Database is at version %d", version)
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
