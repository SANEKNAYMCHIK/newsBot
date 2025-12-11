package repositories

import (
	"context"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type categoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{pool: pool}
}

func (r *categoryRepository) GetAll(ctx context.Context) ([]models.Category, error) {
	query := `SELECT id, name FROM categories ORDER BY name`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *categoryRepository) Create(ctx context.Context, category *models.Category) error {
	query := `
        INSERT INTO categories (name)
        VALUES ($1)
        RETURNING id
    `

	return r.pool.QueryRow(ctx, query, category.Name).Scan(&category.ID)
}
