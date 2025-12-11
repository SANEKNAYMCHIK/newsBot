package services

import (
	"context"
	"errors"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/internal/repositories"
	"github.com/SANEKNAYMCHIK/newsBot/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo   repositories.UserRepository
	jwtManager *auth.JWTManager
}

func NewAuthService(userRepo repositories.UserRepository, jwtManager *auth.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (a *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	existing, _ := a.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashedPasswordStr := string(hashedPassword)
	user := &models.User{
		Email:        &req.Email,
		PasswordHash: &hashedPasswordStr,
		Role:         "user",
	}
	if err := a.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := a.jwtManager.GenerateToken(user)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = nil

	return &models.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (a *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := a.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := a.jwtManager.GenerateToken(user)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = nil

	return &models.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (a *AuthService) RegisterOrLoginTelegram(ctx context.Context, req *models.TelegramRequest) (*models.User, error) {
	user, err := a.userRepo.GetByTelegramID(ctx, req.TgChatID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, errors.New("already use telegram")
	}

	user = &models.User{
		TgChatID:    &req.TgChatID,
		TgFirstName: req.TgFirstName,
		TgUsername:  req.TgUsername,
		Role:        "user",
	}
	if err := a.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
