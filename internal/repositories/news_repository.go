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

type newsRepository struct {
	pool *pgxpool.Pool
}

func NewNewsRepository(pool *pgxpool.Pool) NewsRepository {
	return &newsRepository{pool: pool}
}

func (r *newsRepository) GetByID(ctx context.Context, id int) (*models.NewsItem, error) {
	var news models.NewsItem
	query := `
        SELECT id, title, content, url, published_at, source_id, guid
        FROM news_items 
        WHERE id = $1
    `

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&news.ID,
		&news.Title,
		&news.Content,
		&news.URL,
		&news.PublishedAt,
		&news.SourceID,
		&news.GUID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &news, nil
}

func (r *newsRepository) GetNewsForUser(ctx context.Context, userID int64, page, pageSize int) ([]models.NewsItem, int64, error) {
	countQuery := `
        SELECT COUNT(*) 
        FROM news_items ni
        JOIN user_sources us ON ni.source_id = us.source_id
        WHERE us.user_id = $1
    `

	var total int64
	err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&total)
	log.Println("NewsRepository GetNewsForUser")
	log.Println(total)
	if err != nil {
		return nil, 0, err
	}

	query := `
        SELECT ni.id, ni.title, ni.content, ni.url, ni.published_at, ni.source_id, ni.guid
        FROM news_items ni
        JOIN user_sources us ON ni.source_id = us.source_id
        WHERE us.user_id = $1
        ORDER BY ni.published_at DESC
        LIMIT $2 OFFSET $3
    `

	offset := (page - 1) * pageSize
	log.Print(page, pageSize)
	rows, err := r.pool.Query(ctx, query, userID, pageSize, offset)
	log.Printf("Rows: %v; Err: %v", rows, err)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var news []models.NewsItem
	for rows.Next() {
		var item models.NewsItem
		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Content,
			&item.URL,
			&item.PublishedAt,
			&item.SourceID,
			&item.GUID,
		)
		log.Printf("Item: %v\n", item)
		if err != nil {
			return nil, 0, err
		}
		news = append(news, item)
	}
	log.Println(news)

	return news, total, nil
}

func (r *newsRepository) Create(ctx context.Context, news *models.NewsItem) error {
	query := `
        INSERT INTO news_items 
        (title, content, url, published_at, source_id, guid)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `

	return r.pool.QueryRow(ctx, query,
		news.Title,
		news.Content,
		news.URL,
		news.PublishedAt,
		news.SourceID,
		news.GUID,
	).Scan(&news.ID)
}

func (r *newsRepository) ExistsByGUID(ctx context.Context, sourceID int, guid string) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1 FROM news_items 
            WHERE source_id = $1 AND guid = $2
        )
    `

	var exists bool
	err := r.pool.QueryRow(ctx, query, sourceID, guid).Scan(&exists)
	return exists, err
}

func (r *newsRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM news_items").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count news: %w", err)
	}
	return count, nil
}

func (n *newsRepository) GetBySource(ctx context.Context, sourceID int64, offset, limit int) ([]models.NewsItem, int64, error) {
	query := `
        SELECT 
            id, title, content, url, published_at, source_id,
            COUNT(*) OVER() as total_count
        FROM news_items 
        WHERE source_id = $1
        ORDER BY published_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := n.pool.Query(ctx, query, sourceID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get news by source: %w", err)
	}
	defer rows.Close()

	var newsItems []models.NewsItem
	var totalCount int64

	for rows.Next() {
		var item models.NewsItem
		var total sql.NullInt64
		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Content,
			&item.URL,
			&item.PublishedAt,
			&item.SourceID,
			&total,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan news item: %w", err)
		}
		if total.Valid {
			totalCount = total.Int64
		}
		newsItems = append(newsItems, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return newsItems, totalCount, nil
}

func (n *newsRepository) GetBySourceWithPagination(ctx context.Context, sourceID int64, offset, limit int) ([]models.NewsItem, int64, error) {
	query := `
		SELECT 
			id, title, content, url, published_at, source_id,
			COUNT(*) OVER() as total_count
		FROM news_items 
		WHERE source_id = $1
		ORDER BY published_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := n.pool.Query(ctx, query, sourceID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get news by source with pagination: %w", err)
	}
	defer rows.Close()

	var newsItems []models.NewsItem
	var totalCount int64

	for rows.Next() {
		var item models.NewsItem
		var total sql.NullInt64

		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Content,
			&item.URL,
			&item.PublishedAt,
			&item.SourceID,
			&total,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan news item: %w", err)
		}

		if total.Valid {
			totalCount = total.Int64
		}

		newsItems = append(newsItems, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return newsItems, totalCount, nil
}
