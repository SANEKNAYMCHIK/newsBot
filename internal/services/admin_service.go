package services

import (
	"context"
	"fmt"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/internal/repositories"
)

type AdminService struct {
	userRepo repositories.UserRepository
}

func NewAdminService(userRepo repositories.UserRepository) *AdminService {
	return &AdminService{userRepo: userRepo}
}

func (s *AdminService) MakeAdmin(ctx context.Context, targetUserID, currentUserID int64) error {
	currentUser, err := s.userRepo.GetByID(ctx, currentUserID)
	if err != nil {
		return fmt.Errorf("not found user with ID %d: %w", currentUserID, err)
	}

	if currentUser.Role != "admin" {
		return fmt.Errorf("user should be admin to make another user admin")
	}

	targetUser, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if targetUser.ID == currentUserID {
		return fmt.Errorf("user can't change his own role")
	}

	if targetUser.Role == "admin" {
		return fmt.Errorf("user is already admin")
	}

	targetUser.Role = "admin"
	return s.userRepo.Update(ctx, targetUser)
}

func (s *AdminService) RemoveAdmin(ctx context.Context, targetUserID, currentUserID int64) error {
	currentUser, err := s.userRepo.GetByID(ctx, currentUserID)
	if err != nil {
		return fmt.Errorf("not found user with ID %d: %w", currentUserID, err)
	}

	if currentUser.Role != "admin" {
		return fmt.Errorf("user should be admin to make another user admin")
	}

	targetUser, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if targetUser.ID == currentUserID {
		return fmt.Errorf("user can't change his own role")
	}

	if targetUser.Role != "admin" {
		return fmt.Errorf("selected user is not an admin")
	}

	targetUser.Role = "user"
	return s.userRepo.Update(ctx, targetUser)
}

func (a *AdminService) GetUsers(ctx context.Context, page, pageSize int) (*models.PaginatedResponse[models.User], error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	users, total, err := a.userRepo.GetUsers(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	for i := range users {
		users[i].PasswordHash = nil
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	return &models.PaginatedResponse[models.User]{
		Data:       users,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
