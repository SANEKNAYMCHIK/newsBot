package repositories

import (
	"context"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByTelegramID(ctx context.Context, tgChatID int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	GetUserSubscriptions(ctx context.Context, userID int) ([]int, error)
}
