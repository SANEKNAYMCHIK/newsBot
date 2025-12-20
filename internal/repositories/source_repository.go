package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type sourceRepository struct {
	pool *pgxpool.Pool
}

func NewSourceRepository(pool *pgxpool.Pool) SourceRepository {
	return &sourceRepository{pool: pool}
}

func (r *sourceRepository) GetActive(ctx context.Context) ([]models.Source, error) {
	query := `
        SELECT id, name, url, category_id, is_active
        FROM sources
        WHERE is_active = true
        ORDER BY name
    `

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []models.Source
	for rows.Next() {
		var source models.Source
		err := rows.Scan(
			&source.ID,
			&source.Name,
			&source.URL,
			&source.CategoryID,
			&source.IsActive,
		)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}

	return sources, nil
}

func (r *sourceRepository) GetByID(ctx context.Context, id int) (*models.Source, error) {
	log.Println("GetByID")
	query := `
        SELECT id, name, url, category_id, is_active
        FROM sources
        WHERE id = $1
    `

	var source models.Source
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&source.ID,
		&source.Name,
		&source.URL,
		&source.CategoryID,
		&source.IsActive,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return &source, err
}

func (r *sourceRepository) GetByURL(ctx context.Context, url string) (*models.Source, error) {
	query := `SELECT id, name, url, category_id, is_active,
              FROM sources WHERE url = $1`

	var source models.Source
	err := r.pool.QueryRow(ctx, query, url).Scan(
		&source.ID,
		&source.Name,
		&source.URL,
		&source.CategoryID,
		&source.IsActive,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &source, nil
}

func (r *sourceRepository) Create(ctx context.Context, source *models.Source) error {
	query := `
        INSERT INTO sources (name, url, category_id, is_active)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	return r.pool.QueryRow(ctx, query,
		source.Name,
		source.URL,
		source.CategoryID,
		source.IsActive,
	).Scan(&source.ID)
}

func (r *sourceRepository) Update(ctx context.Context, source *models.Source) error {
	query := `
        UPDATE sources
        SET name = $1, url = $2, category_id = $3, is_active = $4
        WHERE id = $5
    `

	_, err := r.pool.Exec(ctx, query,
		source.Name,
		source.URL,
		source.CategoryID,
		source.IsActive,
		source.ID,
	)

	return err
}

func (r *sourceRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM sources WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *sourceRepository) GetActiveForUser(ctx context.Context, userID int64) ([]models.Source, error) {
	query := `
        SELECT s.id, s.name, s.url, s.category_id, s.is_active
        FROM sources s
        INNER JOIN user_sources us ON s.id = us.source_id
        WHERE us.user_id = $1 AND s.is_active = true
        ORDER BY s.name
    `

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []models.Source
	for rows.Next() {
		var source models.Source
		err := rows.Scan(
			&source.ID,
			&source.Name,
			&source.URL,
			&source.CategoryID,
			&source.IsActive,
		)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}

	return sources, nil
}

func (r *sourceRepository) GetAllWithPagination(ctx context.Context, page, pageSize int) ([]models.Source, int64, error) {
	query := `
        SELECT 
            id, name, url, category_id, is_active,
            COUNT(*) OVER() as total_count
        FROM sources 
        ORDER BY name
        LIMIT $1 OFFSET $2
    `
	offset := (page - 1) * pageSize

	rows, err := r.pool.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get sources with pagination: %w", err)
	}
	defer rows.Close()

	var sources []models.Source
	var totalCount int64

	for rows.Next() {
		var source models.Source
		var total sql.NullInt64

		err := rows.Scan(
			&source.ID,
			&source.Name,
			&source.URL,
			&source.CategoryID,
			&source.IsActive,
			&total,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan source: %w", err)
		}

		if total.Valid {
			totalCount = total.Int64
		}

		sources = append(sources, source)
		log.Println(sources)
		log.Println(totalCount)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return sources, totalCount, nil
}
