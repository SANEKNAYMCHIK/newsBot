package services

import (
	"context"
	"errors"
	"fmt"
	"log"

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

	role := "user"
	if usersCount, err := a.userRepo.Count(ctx); err == nil {
		if usersCount == 0 {
			role = "admin"
			log.Println("The first user automatically set as admin")
		}
	}
	user := &models.User{
		Email:        &req.Email,
		PasswordHash: &hashedPasswordStr,
		Role:         role,
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

func (a *AuthService) RegisterOrUpdateTelegramUser(ctx context.Context, chatID int64, username, firstName string) (*models.User, error) {
	user, err := a.userRepo.GetByTelegramID(ctx, chatID)
	if err == nil && user != nil {
		user.TgChatID = &chatID
		user.TgUsername = &username
		user.TgFirstName = &firstName
		err = a.userRepo.Update(ctx, user)
		if err != nil {
			return user, errors.New("error with updating of new info about user")
		}
		return user, nil
	}

	role := "user"
	if usersCount, err := a.userRepo.Count(ctx); err == nil {
		if usersCount == 0 {
			role = "admin"
			log.Println("The first user automatically set as admin")
		}
	}

	user = &models.User{
		TgChatID:    &chatID,
		TgUsername:  &username,
		TgFirstName: &firstName,
		Role:        role,
	}

	if err := a.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("error of creating user: %w", err)
	}

	log.Printf("New Telegram user: %s (ID: %d, Role: %s)", firstName, chatID, role)
	return user, nil
}
