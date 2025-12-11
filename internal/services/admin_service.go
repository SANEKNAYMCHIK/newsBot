package services

import (
	"context"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/internal/repositories"
)

type AdminService struct {
	userRepo repositories.UserRepository
}

func NewAdminService(userRepo repositories.UserRepository) *AdminService {
	return &AdminService{userRepo: userRepo}
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
