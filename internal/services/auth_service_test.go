package services

import (
	"context"
	"errors"
	"testing"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/pkg/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

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

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Count(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Get(0).(int), args.Error(1)
}

func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := auth.NewJWTManager("test-secret")
	authService := NewAuthService(mockRepo, jwtManager)

	ctx := context.Background()
	req := &models.RegisterRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	mockRepo.On("GetByEmail", ctx, req.Email).Return((*models.User)(nil), nil)
	mockRepo.On("Count", ctx).Return(0, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
		user := args.Get(1).(*models.User)
		user.ID = 1
		user.Role = "admin"
	}).Return(nil)
	response, err := authService.Register(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, response.User)
	assert.Equal(t, int64(1), response.User.ID)
	assert.Equal(t, "admin", response.User.Role)
	assert.NotEmpty(t, response.Token)
	assert.Nil(t, response.User.PasswordHash)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_UserExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := auth.NewJWTManager("test-secret")
	authService := NewAuthService(mockRepo, jwtManager)

	ctx := context.Background()
	req := &models.RegisterRequest{
		Email:    "test@example.com",
		Password: "password",
	}
	existingUser := &models.User{
		ID:    1,
		Email: &req.Email,
		Role:  "user",
	}

	mockRepo.On("GetByEmail", ctx, req.Email).Return(existingUser, nil)
	response, err := authService.Register(ctx, req)

	require.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "user with this email already exists")

	mockRepo.AssertNotCalled(t, "Create")
	mockRepo.AssertExpectations(t)
}

func TestAuthService_LoginSuccess(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := auth.NewJWTManager("test-secret")
	authService := NewAuthService(mockRepo, jwtManager)

	ctx := context.Background()
	req := &models.LoginRequest{
		Email:    "existing@example.com",
		Password: "password",
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	passwordStr := string(hashedPassword)

	existingUser := &models.User{
		ID:           1,
		Email:        &req.Email,
		PasswordHash: &passwordStr,
		Role:         "user",
	}

	mockRepo.On("GetByEmail", ctx, req.Email).Return(existingUser, nil)
	response, err := authService.Login(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, response.User)
	assert.Equal(t, int64(1), response.User.ID)
	assert.NotEmpty(t, response.Token)
	assert.Nil(t, response.User.PasswordHash)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := auth.NewJWTManager("test-secret")
	authService := NewAuthService(mockRepo, jwtManager)

	ctx := context.Background()
	req := &models.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	passwordStr := string(hashedPassword)

	existingUser := &models.User{
		ID:           1,
		Email:        &req.Email,
		PasswordHash: &passwordStr,
		Role:         "user",
	}

	mockRepo.On("GetByEmail", ctx, req.Email).Return(existingUser, nil)
	response, err := authService.Login(ctx, req)

	require.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid credentials")

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := auth.NewJWTManager("test-secret")
	authService := NewAuthService(mockRepo, jwtManager)

	ctx := context.Background()
	req := &models.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	mockRepo.On("GetByEmail", ctx, req.Email).Return((*models.User)(nil), nil)
	response, err := authService.Login(ctx, req)

	require.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid credentials")

	mockRepo.AssertExpectations(t)
}

func TestAuthService_TelegramRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := auth.NewJWTManager("test-secret")
	authService := NewAuthService(mockRepo, jwtManager)

	ctx := context.Background()
	reqChatID := int64(12345)
	userName := "username_test"
	firstName := "test testov"

	mockRepo.On("GetByTelegramID", ctx, reqChatID).Return((*models.User)(nil), nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
		user := args.Get(1).(*models.User)
		user.ID = 1
		user.Role = "admin"
	}).Return(nil)
	mockRepo.On("Count", ctx).Return(0, nil)
	response, err := authService.RegisterOrUpdateTelegramUser(ctx, reqChatID, userName, firstName)

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, response.ID)
	assert.Equal(t, "admin", response.Role)
	assert.Equal(t, int64(1), response.ID)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_TelegramRegister_CreateError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := auth.NewJWTManager("test-secret")
	authService := NewAuthService(mockRepo, jwtManager)

	ctx := context.Background()
	reqChatID := int64(12345)
	userName := "username_test"
	firstName := "test testov"

	mockRepo.On("GetByTelegramID", ctx, reqChatID).Return((*models.User)(nil), nil)
	mockRepo.On("Count", ctx).Return(0, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.User")).Return(errors.New("Error"))
	response, err := authService.RegisterOrUpdateTelegramUser(ctx, reqChatID, userName, firstName)

	require.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "error of creating user")

	mockRepo.AssertExpectations(t)
}

func TestAuthService_TelegramRegister_UserExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtManager := auth.NewJWTManager("test-secret")
	authService := NewAuthService(mockRepo, jwtManager)

	ctx := context.Background()
	reqChatID := int64(12345)
	userName := "username_test"
	firstName := "test testov"

	returnedChatID := &reqChatID
	returnedUserName := &userName
	returnedUser := &models.User{
		ID:         2,
		Role:       "user",
		TgChatID:   returnedChatID,
		TgUsername: returnedUserName,
	}

	mockRepo.On("GetByTelegramID", ctx, reqChatID).Return(returnedUser, nil)
	mockRepo.On("Update", ctx, returnedUser).Return(nil)
	response, err := authService.RegisterOrUpdateTelegramUser(ctx, reqChatID, userName, firstName)

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, int64(2), response.ID)
	assert.Equal(t, "user", response.Role)
	assert.Contains(t, "username_test", *response.TgUsername)

	mockRepo.AssertNotCalled(t, "Count")
	mockRepo.AssertNotCalled(t, "Create")
	mockRepo.AssertExpectations(t)
}
