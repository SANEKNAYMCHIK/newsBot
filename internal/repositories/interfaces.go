package repositories

import (
	"context"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
)

// UserRepository интерфейс для создания
type UserRepository interface {
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByTelegramID(ctx context.Context, tgChatID int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	GetUserSubscriptions(ctx context.Context, userID int) ([]int, error)
}

// NewsRepository интерфейс для работы с новостями
type NewsRepository interface {
	GetByID(ctx context.Context, id int) (*models.NewsItem, error)
	GetLatestNews(ctx context.Context, limit int) ([]models.NewsItem, error)
	GetBySourceID(ctx context.Context, sourceID int, limit int) ([]models.NewsItem, error)
	ExistsByGUID(ctx context.Context, sourceID int, guid string) (bool, error)
	Create(ctx context.Context, news *models.NewsItem) error
	GetUnsentNews(ctx context.Context, userID int, limit int) ([]models.NewsItem, error)
	SearchNews(ctx context.Context, query string, limit int) ([]models.NewsItem, error)
}

// SourceRepository интерфейс для работы с источниками
type SourceRepository interface {
	GetByID(ctx context.Context, id int) (*models.Source, error)
	GetActiveSources(ctx context.Context) ([]models.Source, error)
	GetByURL(ctx context.Context, url string) (*models.Source, error)
	Create(ctx context.Context, source *models.Source) error
	Update(ctx context.Context, source *models.Source) error
	GetSourcesByCategory(ctx context.Context, categoryID int) ([]models.Source, error)
}

// SubscriptionRepository интерфейс для работы с подписками
type SubscriptionRepository interface {
	Subscribe(ctx context.Context, userID, sourceID int) error
	Unsubscribe(ctx context.Context, userID, sourceID int) error
	GetSubscribers(ctx context.Context, sourceID int) ([]int, error)
}

// SentNewsRepository интерфейс для работы с отправленными новостями
type SentNewsRepository interface {
	MarkAsSent(ctx context.Context, newsID, userID int) error
	IsSent(ctx context.Context, newsID, userID int) (bool, error)
	CleanupOldRecords(ctx context.Context, daysOld int) error
}
