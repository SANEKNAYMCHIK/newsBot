package services

import (
	"context"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/internal/repositories"
)

type CategoryService struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func (c *CategoryService) GetCategories(ctx context.Context) ([]models.Category, error) {
	return c.categoryRepo.GetAll(ctx)
}

func (c *CategoryService) CreateCategory(ctx context.Context, name string) (*models.Category, error) {
	category := &models.Category{
		Name: name,
	}

	if err := c.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}
