package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type subscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) *subscriptionRepository {
	return &subscriptionRepository{pool: pool}
}

func (s *subscriptionRepository) GetUserSubscriptions(ctx context.Context, userID int64) ([]models.Source, error) {
	query := `
        SELECT s.id, s.name, s.url, s.category_id, s.is_active
        FROM sources s
        JOIN user_sources us ON s.id = us.source_id
        WHERE us.user_id = $1 AND s.is_active = true
        ORDER BY s.name
    `
	rows, err := s.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user subscriptions: %w", err)
	}
	defer rows.Close()

	var sources []models.Source
	for rows.Next() {
		var source models.Source
		if err := rows.Scan(
			&source.ID,
			&source.Name,
			&source.URL,
			&source.CategoryID,
			&source.IsActive,
		); err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return sources, nil
}

func (s *subscriptionRepository) Subscribe(ctx context.Context, userID, sourceID int64) error {
	query := `
        INSERT INTO user_sources (user_id, source_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id, source_id) DO NOTHING
    `
	res, err := s.pool.Exec(ctx, query, userID, sourceID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("already subscribed to this source")
	}
	return nil
}

func (s *subscriptionRepository) IsSubscribed(ctx context.Context, userID, sourceID int64) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1 FROM user_sources 
            WHERE user_id = $1 AND source_id = $2
        )
    `
	var exists bool
	err := s.pool.QueryRow(ctx, query, userID, sourceID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check subscription: %w", err)
	}
	return exists, nil
}

func (s *subscriptionRepository) Unsubscribe(ctx context.Context, userID, sourceID int64) error {
	query := `DELETE FROM user_sources WHERE user_id = $1 AND source_id = $2`
	res, err := s.pool.Exec(ctx, query, userID, sourceID)
	if err != nil {
		return fmt.Errorf("failed to remove subscription: %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found")
	}
	return nil
}

// func (s *subscriptionRepository) AddSubscription(ctx context.Context, userID, sourceID int64) error {
// 	query := `INSERT INTO user_sources (user_id, source_id) VALUES ($1, $2)`
// 	_, err := s.pool.Exec(ctx, query, userID, sourceID)
// 	if err != nil {
// 		if isDuplicateKeyError(err) {
// 			return fmt.Errorf("already subscribed to this source")
// 		}
// 		return fmt.Errorf("failed to add subscription: %w", err)
// 	}
// 	return nil
// }

// func (s *subscriptionRepository) RemoveSubscription(ctx context.Context, userID, sourceID int64) error {
// 	query := `DELETE FROM user_sources WHERE user_id = $1 AND source_id = $2`
// 	res, err := s.pool.Exec(ctx, query, userID, sourceID)
// 	if err != nil {
// 		return fmt.Errorf("failed to remove subscription: %w", err)
// 	}
// 	if res.RowsAffected() == 0 {
// 		return fmt.Errorf("subscription not found")
// 	}
// 	return nil
// }

// func (s *subscriptionRepository) IsSubscribed(ctx context.Context, userID, sourceID int64) (bool, error) {
// 	query := `SELECT EXISTS(SELECT 1 FROM user_sources WHERE user_id = $1 AND source_id = $2)`
// 	var exists bool
// 	err := s.pool.QueryRow(ctx, query, userID, sourceID).Scan(&exists)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to check subscription: %w", err)
// 	}
// 	return exists, nil
// }

// func isDuplicateKeyError(err error) bool {
// 	return err != nil && (err.Error() == "ERROR: duplicate key value violates unique constraint (SQLSTATE 23505)" ||
// 		err.Error() == "pq: duplicate key value violates unique constraint")
// }

// // Возможно эту или функцию, которая будет просто возвращать всех пользователей, подписанных на этот источник
// // func (r *subscriptionRepository) GetSubscriberCount(ctx context.Context, sourceID int64) (int, error) {
// //     query := `SELECT COUNT(*) FROM user_sources WHERE source_id = $1`

// //     var count int
// //     err := r.pool.QueryRow(ctx, query, sourceID).Scan(&count)
// //     if err != nil {
// //         return 0, fmt.Errorf("failed to get subscriber count: %w", err)
// //     }

// //     return count, nil
// // }
