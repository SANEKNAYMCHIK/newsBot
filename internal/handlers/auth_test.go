// TODO
// Разобраться с тестами для хендлеров
// Так как замокать сервис нельзя, то проводим интеграционное тестирование

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
	"github.com/SANEKNAYMCHIK/newsBot/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Мок UserRepository (как раньше)
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByTelegramID(ctx context.Context, tgChatID int64) (*models.User, error) {
	args := m.Called(ctx, tgChatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUsers(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

func TestAuthHandler_Integration_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	// 1. Создаем мок репозитория
	mockUserRepo := new(MockUserRepository)

	// 2. Создаем реальный JWT менеджер
	jwtManager := auth.NewJWTManager("test-secret")

	// 3. Создаем РЕАЛЬНЫЙ AuthService с моком репозитория
	authService := services.NewAuthService(mockUserRepo, jwtManager)

	// 4. Создаем хендлер с реальным сервисом
	authHandler := NewAuthHandler(authService)

	// 5. Настраиваем мок
	registerReq := models.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Настраиваем поведение мока
	mockUserRepo.On("GetByEmail", mock.Anything, "test@example.com").Return((*models.User)(nil), nil)
	mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	// 6. Настраиваем роутер
	router := gin.Default()
	router.POST("/register", authHandler.Register)

	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.NotNil(t, response.User)
	assert.NotEmpty(t, response.Token)

	mockUserRepo.AssertExpectations(t)
}
