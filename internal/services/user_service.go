package services

import (
	"context"
	"errors"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/internal/repositories"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(
	userRepo repositories.UserRepository,
) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (u *UserService) GetProfile(ctx context.Context, userID int64) (*models.User, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	user.PasswordHash = nil
	return user, nil
}

// func (u *UserService) GetSubscriptions(ctx context.Context, userID int64) ([]*models.Source, error) {
// 	return u.subscriptionRepo.GetUserSubscriptions(ctx, userID)
// }

// func (u *UserService) AddSubscription(ctx context.Context, userID, sourceID int64) error {
// 	return u.subscriptionRepo.AddSubscription(ctx, userID, sourceID)
// }

// func (u *UserService) RemoveSubscription(ctx context.Context, userID, sourceID int64) error {
// 	return u.subscriptionRepo.RemoveSubscription(ctx, userID, sourceID)
// }
