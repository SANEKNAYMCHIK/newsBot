// TODO

package repositories

import (
	"context"
	"testing"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	// Используем testcontainers для запуска PostgreSQL в Docker
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	// Получаем хост и порт контейнера
	host, err := postgresContainer.Host(ctx)
	require.NoError(t, err)

	port, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	connStr := "postgres://testuser:testpass@" + host + ":" + port.Port() + "/testdb?sslmode=disable"

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	// Создаем таблицы
	_, err = pool.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            tg_chat_id BIGINT UNIQUE,
            tg_username VARCHAR(100),
            tg_first_name VARCHAR(100),
            email VARCHAR(255) UNIQUE,
            password_hash VARCHAR(255),
            role VARCHAR(20) DEFAULT 'user' NOT NULL,
            created_at TIMESTAMPTZ DEFAULT NOW()
        );
        
        CREATE TABLE IF NOT EXISTS categories (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) UNIQUE NOT NULL
        );
        
        -- ... остальные таблицы
    `)
	require.NoError(t, err)

	t.Cleanup(func() {
		pool.Close()
		postgresContainer.Terminate(ctx)
	})

	return pool
}

func TestUserRepositoryIntegration_CreateAndGet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Arrange
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)
	ctx := context.Background()

	email := "test@example.com"
	passwordHash := "hashed_password"

	user := &models.User{
		Email:        &email,
		PasswordHash: &passwordHash,
		Role:         "user",
	}

	// Act: создаем пользователя
	err := repo.Create(ctx, user)
	require.NoError(t, err)
	require.NotZero(t, user.ID)

	// Act: получаем пользователя по ID
	fetchedUser, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, fetchedUser)

	// Assert
	assert.Equal(t, user.ID, fetchedUser.ID)
	assert.Equal(t, user.Email, fetchedUser.Email)
	assert.Equal(t, user.PasswordHash, fetchedUser.PasswordHash)
	assert.Equal(t, user.Role, fetchedUser.Role)

	// Act: получаем пользователя по email
	fetchedByEmail, err := repo.GetByEmail(ctx, email)
	require.NoError(t, err)
	require.NotNil(t, fetchedByEmail)
	assert.Equal(t, user.ID, fetchedByEmail.ID)
}
