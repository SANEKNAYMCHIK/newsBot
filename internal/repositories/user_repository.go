package repositories

import (
	"context"
	"fmt"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
        SELECT id, tg_chat_id, tg_username, tg_first_name, email, password_hash, role
        FROM users 
        WHERE id = $1
    `
	var user models.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.TgChatID,
		&user.TgUsername,
		&user.TgFirstName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}

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
	if err != nil {
		return nil, fmt.Errorf("failed to get user by telegram id: %w", err)
	}
	return &user, nil
}

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
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

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
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return err
}

func (r *userRepository) GetUsers(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM users`
	err := r.pool.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
        SELECT id, tg_chat_id, tg_username, tg_first_name, email, password_hash, role
        FROM users
        ORDER BY id
        LIMIT $1 OFFSET $2
    `

	offset := (page - 1) * pageSize
	rows, err := r.pool.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.TgChatID,
			&user.TgUsername,
			&user.TgFirstName,
			&user.Email,
			&user.PasswordHash,
			&user.Role,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *userRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users`
	var count int
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	return count, err
}

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
