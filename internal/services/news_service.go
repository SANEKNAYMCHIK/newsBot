package services

import (
	"context"
	"fmt"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/internal/repositories"
)

type NewsService struct {
	newsRepo         repositories.NewsRepository
	sourceRepo       repositories.SourceRepository
	subscriptionRepo repositories.SubscriptionRepository
}

func NewNewsService(
	newsRepo repositories.NewsRepository,
	sourceRepo repositories.SourceRepository,
	subscriptionRepo repositories.SubscriptionRepository,
) *NewsService {
	return &NewsService{
		newsRepo:         newsRepo,
		sourceRepo:       sourceRepo,
		subscriptionRepo: subscriptionRepo,
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

	totalPages := 0
	if total > 0 {
		totalPages = int(total) / pageSize
		if int(total)%pageSize > 0 {
			totalPages++
		}
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

func (s *NewsService) GetNewsBySource(
	ctx context.Context,
	sourceID int64,
	userID int64,
	page, pageSize int,
) (*models.PaginatedResponse[models.NewsResponse], error) {
	source, err := s.sourceRepo.GetByID(ctx, int(sourceID))
	if err != nil || source == nil {
		return nil, fmt.Errorf("source doesn't exist")
	}

	subscribed, err := s.subscriptionRepo.IsSubscribed(ctx, userID, sourceID)
	if err != nil || !subscribed {
		return nil, fmt.Errorf("user doesn't subscribe on the source")
	}

	offset := (page - 1) * pageSize
	newsItems, total, err := s.newsRepo.GetBySourceWithPagination(ctx, sourceID, offset, pageSize)
	if err != nil {
		return nil, err
	}

	var data []models.NewsResponse
	for _, item := range newsItems {
		var content *string
		if item.Content != nil {
			content = item.Content
		}
		data = append(data, models.NewsResponse{
			ID:          item.ID,
			Title:       item.Title,
			Content:     content,
			URL:         item.URL,
			PublishedAt: item.PublishedAt,
			SourceID:    item.SourceID,
			SourceName:  source.Name,
		})
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(total) / pageSize
		if int(total)%pageSize > 0 {
			totalPages++
		}
	}

	return &models.PaginatedResponse[models.NewsResponse]{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
