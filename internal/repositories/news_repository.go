package repositories

import (
	"context"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type newsRepository struct {
	pool *pgxpool.Pool
}

func NewNewsRepository(p *pgxpool.Pool) NewsRepository {
	return &newsRepository{pool: p}
}

// GetByID возвращает новость по ID
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

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &news, nil
}

// GetLatestNews возвращает последние новости
func (r *newsRepository) GetLatestNews(ctx context.Context, limit int) ([]models.NewsItem, error) {
	query := `
        SELECT id, title, content, url, published_at, source_id, guid
        FROM news_items 
        ORDER BY published_at DESC 
        LIMIT $1
    `

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
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
		if err != nil {
			return nil, err
		}
		news = append(news, item)
	}

	return news, nil
}

// GetBySourceID возвращает новости по источнику
func (r *newsRepository) GetBySourceID(ctx context.Context, sourceID int, limit int) ([]models.NewsItem, error) {
	query := `
        SELECT id, title, content, url, published_at, source_id, guid
        FROM news_items 
        WHERE source_id = $1
        ORDER BY published_at DESC 
        LIMIT $2
    `

	rows, err := r.pool.Query(ctx, query, sourceID, limit)
	if err != nil {
		return nil, err
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
		if err != nil {
			return nil, err
		}
		news = append(news, item)
	}

	return news, nil
}

// ExistsByGUID проверяет существует ли новость с таким GUID
func (r *newsRepository) ExistsByGUID(ctx context.Context, sourceID int, guid string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM news_items WHERE source_id = $1 AND guid = $2)`

	err := r.pool.QueryRow(ctx, query, sourceID, guid).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// Create создает новую новость
func (r *newsRepository) Create(ctx context.Context, news *models.NewsItem) error {
	query := `
        INSERT INTO news_items (title, content, url, published_at, source_id, guid)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `

	err := r.pool.QueryRow(ctx, query,
		news.Title,
		news.Content,
		news.URL,
		news.PublishedAt,
		news.SourceID,
		news.GUID,
	).Scan(&news.ID)

	return err
}

// GetUnsentNews возвращает непрочитанные новости для пользователя
func (r *newsRepository) GetUnsentNews(ctx context.Context, userID int, limit int) ([]models.NewsItem, error) {
	query := `
        SELECT ni.id, ni.title, ni.content, ni.url, ni.published_at, ni.source_id, ni.guid
        FROM news_items ni
        WHERE ni.source_id IN (
            SELECT source_id FROM user_sources WHERE user_id = $1
        )
        AND NOT EXISTS (
            SELECT 1 FROM sent_news sn 
            WHERE sn.news_id = ni.id AND sn.user_id = $1
        )
        ORDER BY ni.published_at DESC
        LIMIT $2
    `

	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
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
		if err != nil {
			return nil, err
		}
		news = append(news, item)
	}

	return news, nil
}

// SearchNews ищет новости по заголовку и содержанию
func (r *newsRepository) SearchNews(ctx context.Context, query string, limit int) ([]models.NewsItem, error) {
	searchQuery := `
        SELECT id, title, content, url, published_at, source_id, guid
        FROM news_items 
        WHERE title ILIKE $1 OR content ILIKE $1
        ORDER BY published_at DESC
        LIMIT $2
    `

	rows, err := r.pool.Query(ctx, searchQuery, "%"+query+"%", limit)
	if err != nil {
		return nil, err
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
		if err != nil {
			return nil, err
		}
		news = append(news, item)
	}

	return news, nil
}
