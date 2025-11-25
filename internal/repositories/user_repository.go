package repositories

import (
	"context"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(p *pgxpool.Pool) UserRepository {
	return &userRepository{pool: p}
}

// GetByID возвращает пользователя по ID
func (r *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	query := `
        SELECT id, tg_chat_id, tg_username, tg_first_name, email, password_hash, role
        FROM users 
        WHERE id = $1
    `

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.TgChatID,
		&user.TgUsername,
		&user.TgFirstName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByTelegramID возвращает пользователя по Telegram Chat ID
func (r *userRepository) GetByTelegramID(ctx context.Context, tgChatID int64) (*models.User, error) {
	var user models.User
	query := `
        SELECT id, tg_chat_id, tg_username, tg_first_name, email, password_hash, role
        FROM users 
        WHERE tg_chat_id = $1
    `

	err := r.pool.QueryRow(ctx, query, tgChatID).Scan(
		&user.ID,
		&user.TgChatID,
		&user.TgUsername,
		&user.TgFirstName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail возвращает пользователя по email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `
        SELECT id, tg_chat_id, tg_username, tg_first_name, email, password_hash, role
        FROM users 
        WHERE email = $1
    `

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.TgChatID,
		&user.TgUsername,
		&user.TgFirstName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Create создает нового пользователя
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (tg_chat_id, tg_username, tg_first_name, email, password_hash, role)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `

	err := r.pool.QueryRow(ctx, query,
		user.TgChatID,
		user.TgUsername,
		user.TgFirstName,
		user.Email,
		user.PasswordHash,
		user.Role,
	).Scan(&user.ID)

	return err
}

// Update обновляет данные пользователя
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
        UPDATE users 
        SET tg_chat_id = $1, tg_username = $2, tg_first_name = $3, 
            email = $4, password_hash = $5, role = $6
        WHERE id = $7
    `

	_, err := r.pool.Exec(ctx, query,
		user.TgChatID,
		user.TgUsername,
		user.TgFirstName,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.ID,
	)

	return err
}

// GetUserSubscriptions возвращает подписки пользователя
func (r *userRepository) GetUserSubscriptions(ctx context.Context, userID int) ([]int, error) {
	query := `
        SELECT source_id 
        FROM user_sources 
        WHERE user_id = $1
    `

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sourceIDs []int
	for rows.Next() {
		var sourceID int
		if err := rows.Scan(&sourceID); err != nil {
			return nil, err
		}
		sourceIDs = append(sourceIDs, sourceID)
	}

	return sourceIDs, nil
}
