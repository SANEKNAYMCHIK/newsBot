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
	Count(ctx context.Context) (int, error)
	Update(ctx context.Context, user *models.User) error
}

type SubscriptionRepository interface {
	GetUserSubscriptions(ctx context.Context, userID int64) ([]models.Source, error)
	Subscribe(ctx context.Context, userID, sourceID int64) error
	IsSubscribed(ctx context.Context, userID, sourceID int64) (bool, error)
	Unsubscribe(ctx context.Context, userID, sourceID int64) error
}

type NewsRepository interface {
	GetNewsForUser(ctx context.Context, userID int64, page, pageSize int) ([]models.NewsItem, int64, error)
	GetByID(ctx context.Context, id int) (*models.NewsItem, error)
	GetBySource(ctx context.Context, sourceID int64, offset, limit int) ([]models.NewsItem, int64, error)
	GetBySourceWithPagination(ctx context.Context, sourceID int64, offset, limit int) ([]models.NewsItem, int64, error)
	ExistsByGUID(ctx context.Context, sourceID int, guid string) (bool, error)
	Create(ctx context.Context, news *models.NewsItem) error
	Count(ctx context.Context) (int64, error)
}

type SourceRepository interface {
	GetActive(ctx context.Context) ([]models.Source, error)
	GetByID(ctx context.Context, id int) (*models.Source, error)
	GetByURL(ctx context.Context, url string) (*models.Source, error)
	Create(ctx context.Context, source *models.Source) error
	Update(ctx context.Context, source *models.Source) error
	Delete(ctx context.Context, id int) error
	GetActiveForUser(ctx context.Context, userID int64) ([]models.Source, error)
	GetAllWithPagination(ctx context.Context, offset, limit int) ([]models.Source, int64, error)
}

type CategoryRepository interface {
	GetAll(ctx context.Context) ([]models.Category, error)
	Create(ctx context.Context, category *models.Category) error
}
