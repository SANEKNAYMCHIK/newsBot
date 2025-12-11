package services

import (
	"context"
	"errors"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/internal/repositories"
)

type SourceService struct {
	sourceRepo repositories.SourceRepository
}

func NewSourceService(sourceRepo repositories.SourceRepository) *SourceService {
	return &SourceService{sourceRepo: sourceRepo}
}

func (s *SourceService) GetActiveSources(ctx context.Context) ([]models.Source, error) {
	return s.sourceRepo.GetActive(ctx)
}

func (s *SourceService) CreateSource(ctx context.Context, req *models.CreateSourceRequest) (*models.Source, error) {
	source := &models.Source{
		Name:       req.Name,
		URL:        req.URL,
		CategoryID: req.CategoryID,
		IsActive:   req.IsActive,
	}

	if err := s.sourceRepo.Create(ctx, source); err != nil {
		return nil, err
	}

	return source, nil
}

func (s *SourceService) UpdateSource(ctx context.Context, sourceID int, req *models.UpdateSourceRequest) (*models.Source, error) {
	source, err := s.sourceRepo.GetByID(ctx, sourceID)
	if err != nil {
		return nil, err
	}
	if source == nil {
		return nil, errors.New("source not found")
	}

	if req.Name != nil {
		source.Name = *req.Name
	}
	if req.URL != nil {
		source.URL = *req.URL
	}
	if req.CategoryID != nil {
		source.CategoryID = req.CategoryID
	}
	if req.IsActive != nil {
		source.IsActive = *req.IsActive
	}

	if err := s.sourceRepo.Update(ctx, source); err != nil {
		return nil, err
	}

	return source, nil
}

func (s *SourceService) DeleteSource(ctx context.Context, sourceID int) error {
	return s.sourceRepo.Delete(ctx, sourceID)
}
