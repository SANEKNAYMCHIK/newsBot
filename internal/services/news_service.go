package services

import (
	"context"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/internal/repositories"
)

type NewsService struct {
	newsRepo   repositories.NewsRepository
	sourceRepo repositories.SourceRepository
}

func NewNewsService(
	newsRepo repositories.NewsRepository,
	sourceRepo repositories.SourceRepository,
) *NewsService {
	return &NewsService{
		newsRepo:   newsRepo,
		sourceRepo: sourceRepo,
	}
}

func (n *NewsService) GetNews(ctx context.Context, userID int64, page, pageSize int) (*models.PaginatedResponse[models.NewsResponse], error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	news, total, err := n.newsRepo.GetNewsForUser(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	var response []models.NewsResponse
	for _, item := range news {
		source, err := n.sourceRepo.GetByID(ctx, int(item.SourceID))
		if err != nil {
			continue
		}
		response = append(response, models.NewsResponse{
			ID:          item.ID,
			Title:       item.Title,
			Content:     item.Content,
			URL:         item.URL,
			PublishedAt: item.PublishedAt,
			SourceID:    item.SourceID,
			SourceName:  source.Name,
			CategoryID:  source.CategoryID,
		})
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &models.PaginatedResponse[models.NewsResponse]{
		Data:       response,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *NewsService) GetNewsByID(ctx context.Context, newsID int) (*models.NewsResponse, error) {
	news, err := s.newsRepo.GetByID(ctx, newsID)
	if err != nil {
		return nil, err
	}
	if news == nil {
		return nil, nil
	}

	source, err := s.sourceRepo.GetByID(ctx, int(news.SourceID))
	if err != nil {
		return nil, err
	}

	return &models.NewsResponse{
		ID:          news.ID,
		Title:       news.Title,
		Content:     news.Content,
		URL:         news.URL,
		PublishedAt: news.PublishedAt,
		SourceID:    news.SourceID,
		SourceName:  source.Name,
		CategoryID:  source.CategoryID,
	}, nil
}
