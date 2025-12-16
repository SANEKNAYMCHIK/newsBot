package services

import (
	"context"
	"testing"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *models.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetAll(ctx context.Context) ([]models.Category, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id int) (*models.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(ctx context.Context, category *models.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCategoryService_GetCategories_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	categoryService := NewCategoryService(mockRepo)
	ctx := context.Background()

	testCategories := []models.Category{
		{
			ID:   1,
			Name: "Technology",
		},
		{
			ID:   2,
			Name: "Sports",
		},
	}

	mockRepo.On("GetAll", ctx).Return(testCategories, nil)

	categories, err := categoryService.GetCategories(ctx)

	require.NoError(t, err)
	assert.NotNil(t, categories)
	assert.Len(t, categories, 2)
	assert.Equal(t, "Technology", categories[0].Name)
	assert.Equal(t, "Sports", categories[1].Name)

	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetCategories_Empty(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	categoryService := NewCategoryService(mockRepo)

	ctx := context.Background()

	mockRepo.On("GetAll", ctx).Return([]models.Category{}, nil)

	categories, err := categoryService.GetCategories(ctx)

	require.NoError(t, err)
	assert.NotNil(t, categories)
	assert.Len(t, categories, 0)

	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetCategories_RepositoryError(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	categoryService := NewCategoryService(mockRepo)

	ctx := context.Background()

	mockRepo.On("GetAll", ctx).Return(([]models.Category)(nil), assert.AnError)

	categories, err := categoryService.GetCategories(ctx)

	require.Error(t, err)
	assert.Nil(t, categories)

	mockRepo.AssertExpectations(t)
}

func TestCategoryService_CreateCategory_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	categoryService := NewCategoryService(mockRepo)

	ctx := context.Background()
	categoryName := "New Category"

	mockRepo.On("Create", ctx, mock.MatchedBy(func(c *models.Category) bool {
		return c.Name == categoryName
	})).Run(func(args mock.Arguments) {
		category := args.Get(1).(*models.Category)
		category.ID = 1
	}).Return(nil)

	category, err := categoryService.CreateCategory(ctx, categoryName)

	require.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, int64(1), category.ID)
	assert.Equal(t, categoryName, category.Name)

	mockRepo.AssertExpectations(t)
}

func TestCategoryService_CreateCategory_RepositoryError(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	categoryService := NewCategoryService(mockRepo)

	ctx := context.Background()
	categoryName := "New Category"

	mockRepo.On("Create", ctx, mock.MatchedBy(func(c *models.Category) bool {
		return c.Name == categoryName
	})).Return(assert.AnError)

	category, err := categoryService.CreateCategory(ctx, categoryName)

	require.Error(t, err)
	assert.Nil(t, category)

	mockRepo.AssertExpectations(t)
}
