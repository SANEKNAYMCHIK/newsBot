package services

import (
	"context"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/pkg/auth"
)

type AuthService struct {
	jwtManager *auth.JWTManager
	// userRepo
}

func NewAuthService(jwtManager *auth.JWTManager) *AuthService {
	return &AuthService{
		jwtManager: jwtManager,
	}
}

func (a *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, string, error) {
	// TODO: Check existing password of user
	// TODO: Hashing password
	// TODO: Save into Database

	// TODO: Creating user object
	user := &models.User{}

	token, err := a.jwtManager.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (a *AuthService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	// TODO: Check in Database
	// TODO: Check password
	return &models.User{}, "", nil
}
