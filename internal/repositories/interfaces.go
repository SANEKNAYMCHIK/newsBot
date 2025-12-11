package repositories

import (
	"context"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByTelegramID(ctx context.Context, tgChatID int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	GetUsers(ctx context.Context, page, pageSize int) ([]models.User, int64, error)
	// Update(ctx context.Context, user *models.User) error
	// GetUserSubscriptions(ctx context.Context, userID int) ([]int, error)
}

type SubscriptionRepository interface {
	GetUserSubscriptions(ctx context.Context, userID int64) ([]models.Source, error)
	Subscribe(ctx context.Context, userID, sourceID int64) error
	IsSubscribed(ctx context.Context, userID, sourceID int64) (bool, error)
	Unsubscribe(ctx context.Context, userID, sourceID int64) error
	// AddSubscription(ctx context.Context, userID, sourceID int64) error
	// RemoveSubscription(ctx context.Context, userID, sourceID int64) error
	// GetSubscriberCount(ctx context.Context, sourceID int64) (int, error)
}

type NewsRepository interface {
	GetNewsForUser(ctx context.Context, userID int64, page, pageSize int) ([]models.NewsItem, int64, error)
	GetByID(ctx context.Context, id int) (*models.NewsItem, error)
	// GetByID(ctx context.Context, id int) (*models.NewsItem, error)
	// GetLatestNews(ctx context.Context, limit int) ([]models.NewsItem, error)
	// GetBySourceID(ctx context.Context, sourceID int, limit int) ([]models.NewsItem, error)
	// ExistsByGUID(ctx context.Context, sourceID int, guid string) (bool, error)
	// Create(ctx context.Context, news *models.NewsItem) error
	// GetUnsentNews(ctx context.Context, userID int, limit int) ([]models.NewsItem, error)
	// SearchNews(ctx context.Context, query string, limit int) ([]models.NewsItem, error)
}

type SourceRepository interface {
	GetActive(ctx context.Context) ([]models.Source, error)
	GetByID(ctx context.Context, id int) (*models.Source, error)
	Create(ctx context.Context, source *models.Source) error
	Update(ctx context.Context, source *models.Source) error
	Delete(ctx context.Context, id int) error
	// GetByURL(ctx context.Context, url string) (*models.Source, error)
	// GetSourcesByCategory(ctx context.Context, categoryID int64) ([]models.Source, error)
}

type CategoryRepository interface {
	GetAll(ctx context.Context) ([]models.Category, error)
	Create(ctx context.Context, category *models.Category) error
	// MarkAsSent(ctx context.Context, newsID, userID int) error
	// IsSent(ctx context.Context, newsID, userID int) (bool, error)
	// CleanupOldRecords(ctx context.Context, daysOld int) error
}
