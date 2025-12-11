package services

import (
	"context"
	"errors"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/internal/repositories"
)

type SubscriptionService struct {
	subscriptionRepo repositories.SubscriptionRepository
	sourceRepo       repositories.SourceRepository
}

func NewSubscriptionService(
	subscriptionRepo repositories.SubscriptionRepository,
	sourceRepo repositories.SourceRepository,
) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
		sourceRepo:       sourceRepo,
	}
}

func (s *SubscriptionService) GetUserSubscriptions(ctx context.Context, userID int64) ([]models.SubscriptionResponse, error) {
	sources, err := s.subscriptionRepo.GetUserSubscriptions(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []models.SubscriptionResponse
	for _, source := range sources {
		result = append(result, models.SubscriptionResponse{
			SourceID:   source.ID,
			SourceName: source.Name,
			CategoryID: source.CategoryID,
			IsActive:   source.IsActive,
		})
	}

	return result, nil
}

func (s *SubscriptionService) AddSubscription(ctx context.Context, userID, sourceID int64) error {
	source, err := s.sourceRepo.GetByID(ctx, int(sourceID))
	if err != nil {
		return err
	}
	if source == nil {
		return errors.New("source not found")
	}
	if !source.IsActive {
		return errors.New("source is not active")
	}

	subscribed, err := s.subscriptionRepo.IsSubscribed(ctx, userID, sourceID)
	if err != nil {
		return err
	}
	if subscribed {
		return errors.New("already subscribed")
	}

	return s.subscriptionRepo.Subscribe(ctx, userID, sourceID)
}

func (s *SubscriptionService) RemoveSubscription(ctx context.Context, userID, sourceID int64) error {
	subscribed, err := s.subscriptionRepo.IsSubscribed(ctx, userID, sourceID)
	if err != nil {
		return err
	}
	if !subscribed {
		return errors.New("not subscribed")
	}

	return s.subscriptionRepo.Unsubscribe(ctx, userID, sourceID)
}
