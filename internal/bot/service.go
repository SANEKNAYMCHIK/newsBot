package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/internal/repositories"
	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
)

type BotService struct {
	authService      *services.AuthService
	userRepo         repositories.UserRepository
	sourceRepo       repositories.SourceRepository
	newsRepo         repositories.NewsRepository
	subscriptionRepo repositories.SubscriptionRepository
	categoryRepo     repositories.CategoryRepository
	newsService      *services.NewsService
	adminService     *services.AdminService
	categoryService  *services.CategoryService
	sourceService    *services.SourceService
	refreshService   *services.RefreshService
}

type NewsWithSource struct {
	ID          int64
	Title       string
	Content     string
	URL         string
	PublishedAt time.Time
	SourceID    int64
	SourceName  string
}

func NewBotService(
	authService *services.AuthService,
	sourceRepo repositories.SourceRepository,
	userRepo repositories.UserRepository,
	newsRepo repositories.NewsRepository,
	subscriptionRepo repositories.SubscriptionRepository,
	categoryRepo repositories.CategoryRepository,
	newsService *services.NewsService,
	adminService *services.AdminService,
	categoryService *services.CategoryService,
	sourceService *services.SourceService,
	refreshService *services.RefreshService,
) *BotService {
	return &BotService{
		authService:      authService,
		sourceRepo:       sourceRepo,
		userRepo:         userRepo,
		newsRepo:         newsRepo,
		subscriptionRepo: subscriptionRepo,
		categoryRepo:     categoryRepo,
		newsService:      newsService,
		adminService:     adminService,
		categoryService:  categoryService,
		sourceService:    sourceService,
		refreshService:   refreshService,
	}
}

func (s *BotService) GetAllSources(ctx context.Context, page, pageSize int) (*models.PaginatedResponse[models.Source], error) {
	offset := (page - 1) * pageSize
	sources, total, err := s.sourceRepo.GetAllWithPagination(ctx, offset, pageSize)
	if err != nil {
		return nil, err
	}

	var activeSources []models.Source
	for _, source := range sources {
		if source.IsActive {
			activeSources = append(activeSources, source)
		}
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(total) / pageSize
		if int(total)%pageSize > 0 {
			totalPages++
		}
	}

	return &models.PaginatedResponse[models.Source]{
		Data:       activeSources,
		Total:      int64(len(activeSources)),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *BotService) GetNewsForUser(ctx context.Context, userID int64, limit int) ([]NewsWithSource, error) {
	return s.GetNewsForUserLegacy(ctx, userID, limit)
}

func (s *BotService) GetNewsForUserLegacy(ctx context.Context, userID int64, limit int) ([]NewsWithSource, error) {
	newsItems, _, err := s.newsRepo.GetNewsForUser(ctx, userID, 0, limit)
	if err != nil {
		return nil, err
	}

	var result []NewsWithSource
	for _, item := range newsItems {
		source, err := s.sourceRepo.GetByID(ctx, int(item.SourceID))
		if err != nil {
			continue
		}

		content := ""
		if item.Content != nil {
			content = *item.Content
		}

		result = append(result, NewsWithSource{
			ID:          item.ID,
			Title:       item.Title,
			Content:     content,
			URL:         item.URL,
			PublishedAt: item.PublishedAt,
			SourceID:    item.SourceID,
			SourceName:  source.Name,
		})
	}

	return result, nil
}

func (s *BotService) GetNewsForUserWithPagination(ctx context.Context, userID int64, page, pageSize int) (*models.PaginatedResponse[NewsWithSource], error) {
	log.Println("GetNewsForUser")
	// offset := (page - 1) * pageSize
	newsItems, total, err := s.newsRepo.GetNewsForUser(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	var data []NewsWithSource
	log.Printf("newsItems Size: %d\n", total)
	log.Println(newsItems)
	for _, item := range newsItems {
		source, err := s.sourceRepo.GetByID(ctx, int(item.SourceID))
		if err != nil {
			continue
		}
		content := ""
		if item.Content != nil {
			content = *item.Content
		}
		data = append(data, NewsWithSource{
			ID:          item.ID,
			Title:       item.Title,
			Content:     content,
			URL:         item.URL,
			PublishedAt: item.PublishedAt,
			SourceID:    item.SourceID,
			SourceName:  source.Name,
		})
		log.Println(data)
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(total) / pageSize
		if int(total)%pageSize > 0 {
			totalPages++
		}
	}

	return &models.PaginatedResponse[NewsWithSource]{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// func (s *BotService) GetNewsBySource(ctx context.Context, sourceID, userID int64, limit int) ([]NewsWithSource, error) {
// 	return s.GetNewsBySourceLegacy(ctx, sourceID, userID, limit)
// }

// func (s *BotService) GetNewsBySourceLegacy(ctx context.Context, sourceID, userID int64, limit int) ([]NewsWithSource, error) {
// 	subscribed, err := s.subscriptionRepo.IsSubscribed(ctx, userID, sourceID)
// 	if err != nil || !subscribed {
// 		return nil, fmt.Errorf("вы не подписаны на этот источник")
// 	}
// 	newsItems, _, err := s.newsRepo.GetBySource(ctx, sourceID, 5, limit)
// 	if err != nil {
// 		return nil, err
// 	}
// 	source, err := s.sourceRepo.GetByID(ctx, int(sourceID))
// 	if err != nil {
// 		return nil, err
// 	}

// 	var result []NewsWithSource
// 	for _, item := range newsItems {
// 		content := ""
// 		if item.Content != nil {
// 			content = *item.Content
// 		}
// 		result = append(result, NewsWithSource{
// 			ID:          item.ID,
// 			Title:       item.Title,
// 			Content:     content,
// 			URL:         item.URL,
// 			PublishedAt: item.PublishedAt,
// 			SourceID:    item.SourceID,
// 			SourceName:  source.Name,
// 		})
// 	}
// 	return result, nil
// }

func (s *BotService) GetNewsBySourceWithPagination(ctx context.Context, sourceID, userID int64, page, pageSize int) (*models.PaginatedResponse[NewsWithSource], error) {
	subscribed, err := s.subscriptionRepo.IsSubscribed(ctx, userID, sourceID)
	if err != nil || !subscribed {
		return nil, fmt.Errorf("вы не подписаны на этот источник")
	}
	offset := (page - 1) * pageSize
	newsItems, total, err := s.newsRepo.GetBySourceWithPagination(ctx, sourceID, offset, pageSize)
	if err != nil {
		return nil, err
	}

	source, err := s.sourceRepo.GetByID(ctx, int(sourceID))
	if err != nil {
		return nil, err
	}

	var data []NewsWithSource
	for _, item := range newsItems {
		content := ""
		if item.Content != nil {
			content = *item.Content
		}

		data = append(data, NewsWithSource{
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

	return &models.PaginatedResponse[NewsWithSource]{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *BotService) GetUserSubscriptions(ctx context.Context, userID int64) ([]models.Source, error) {
	subscriptions, err := s.subscriptionRepo.GetUserSubscriptions(ctx, userID)
	if err != nil {
		return nil, err
	}

	var sources []models.Source
	for _, sub := range subscriptions {
		log.Println(sub)
		log.Println(sub.ID, int(sub.ID), sub.ID)
		source, err := s.sourceRepo.GetByID(ctx, int(sub.ID))
		log.Println(source)
		log.Println(err)
		if err != nil {
			log.Printf("Error getting source %d: %v", sub.ID, err)
			continue
		}
		sources = append(sources, *source)
	}

	return sources, nil
}

func (s *BotService) GetAllActiveSources(ctx context.Context) ([]models.Source, error) {
	return s.sourceRepo.GetActive(ctx)
}

func (s *BotService) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	return s.categoryRepo.GetAll(ctx)
}

func (s *BotService) SubscribeUser(ctx context.Context, userID int64, sourceID int) error {
	exists, err := s.subscriptionRepo.IsSubscribed(ctx, userID, int64(sourceID))
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("already subscribed")
	}

	return s.subscriptionRepo.Subscribe(ctx, userID, int64(sourceID))
}

func (s *BotService) UnsubscribeUser(ctx context.Context, userID int64, sourceID int) error {
	return s.subscriptionRepo.Unsubscribe(ctx, userID, int64(sourceID))
}

func (s *BotService) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}
	return user.Role == "admin", nil
}

func (s *BotService) AddSource(ctx context.Context, name, url string, categoryID, userID int64) error {
	existing, err := s.sourceRepo.GetByURL(ctx, url)
	if err == nil && existing != nil {
		return fmt.Errorf("источник с таким URL уже существует")
	}
	catID := &categoryID
	source := &models.Source{
		Name:       name,
		URL:        url,
		CategoryID: catID,
		IsActive:   true,
	}

	err = s.sourceRepo.Create(ctx, source)
	return err
}

func (s *BotService) RequestNewsUpdate(ctx context.Context, userID int64) (string, error) {
	req, err := s.refreshService.RequestRefresh(ctx, userID)
	if err != nil {
		return "", err
	}
	return req.ID, nil
}

func (s *BotService) GetUpdateStatus(ctx context.Context, requestID string) (*services.RefreshRequest, bool) {
	return s.refreshService.GetRequestStatus(requestID)
}

func (s *BotService) GetUsers(ctx context.Context, page, pageSize int) (*models.PaginatedResponse[models.User], error) {
	return s.adminService.GetUsers(ctx, page, pageSize)
}

func (s *BotService) MakeAdmin(ctx context.Context, targetUserID, currentUserID int64) error {
	return s.adminService.MakeAdmin(ctx, targetUserID, currentUserID)
}

func (s *BotService) RemoveAdmin(ctx context.Context, targetUserID, currentUserID int64) error {
	return s.adminService.RemoveAdmin(ctx, targetUserID, currentUserID)
}

func (s *BotService) CreateCategory(ctx context.Context, name string) (*models.Category, error) {
	return s.categoryService.CreateCategory(ctx, name)
}

func (s *BotService) UpdateSource(ctx context.Context, sourceID int, isActive bool) error {
	req := &models.UpdateSourceRequest{
		IsActive: &isActive,
	}
	_, err := s.sourceService.UpdateSource(ctx, sourceID, req)
	return err
}

func (s *BotService) DeleteSource(ctx context.Context, sourceID int) error {
	return s.sourceService.DeleteSource(ctx, sourceID)
}

func (s *BotService) GetSystemStats(ctx context.Context) (map[string]interface{}, error) {
	users, _, err := s.userRepo.GetUsers(ctx, 1, 1)
	if err != nil {
		return nil, err
	}

	sources, err := s.sourceRepo.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	newsCount, err := s.newsRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"users_count":   len(users),
		"sources_count": len(sources),
		"news_count":    newsCount,
	}

	return stats, nil
}
