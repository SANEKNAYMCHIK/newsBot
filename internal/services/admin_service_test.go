package services

import (
	"context"
	"testing"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAdminService_MakeAdmin_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	adminService := NewAdminService(mockRepo)

	ctx := context.Background()
	currentUserID := int64(1)
	targetUserID := int64(2)

	currentUser := &models.User{
		ID:   currentUserID,
		Role: "admin",
	}

	targetUser := &models.User{
		ID:   targetUserID,
		Role: "user",
	}

	mockRepo.On("GetByID", ctx, currentUserID).Return(currentUser, nil)
	mockRepo.On("GetByID", ctx, targetUserID).Return(targetUser, nil)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(u *models.User) bool {
		return u.ID == targetUserID && u.Role == "admin"
	})).Return(nil)

	err := adminService.MakeAdmin(ctx, targetUserID, currentUserID)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAdminService_MakeAdmin_CurrentUserNotAdmin(t *testing.T) {
	mockRepo := new(MockUserRepository)
	adminService := NewAdminService(mockRepo)

	ctx := context.Background()
	currentUserID := int64(1)
	targetUserID := int64(2)

	currentUser := &models.User{
		ID:   currentUserID,
		Role: "user",
	}

	mockRepo.On("GetByID", ctx, currentUserID).Return(currentUser, nil)

	err := adminService.MakeAdmin(ctx, targetUserID, currentUserID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "user should be admin")
	mockRepo.AssertNotCalled(t, "Update")
	mockRepo.AssertExpectations(t)
}

func TestAdminService_MakeAdmin_TargetUserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	adminService := NewAdminService(mockRepo)

	ctx := context.Background()
	currentUserID := int64(1)
	targetUserID := int64(999)

	currentUser := &models.User{
		ID:   currentUserID,
		Role: "admin",
	}

	mockRepo.On("GetByID", ctx, currentUserID).Return(currentUser, nil)
	mockRepo.On("GetByID", ctx, targetUserID).Return((*models.User)(nil), assert.AnError)

	err := adminService.MakeAdmin(ctx, targetUserID, currentUserID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	mockRepo.AssertNotCalled(t, "Update")
	mockRepo.AssertExpectations(t)
}

func TestAdminService_MakeAdmin_ChangeOwnRole(t *testing.T) {
	mockRepo := new(MockUserRepository)
	adminService := NewAdminService(mockRepo)

	ctx := context.Background()
	userID := int64(1)

	currentUser := &models.User{
		ID:   userID,
		Role: "admin",
	}

	mockRepo.On("GetByID", ctx, userID).Return(currentUser, nil)

	err := adminService.MakeAdmin(ctx, userID, userID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "can't change his own role")
	mockRepo.AssertNotCalled(t, "Update")
	mockRepo.AssertExpectations(t)
}

func TestAdminService_MakeAdmin_UserAlreadyAdmin(t *testing.T) {
	mockRepo := new(MockUserRepository)
	adminService := NewAdminService(mockRepo)

	ctx := context.Background()
	currentUserID := int64(1)
	targetUserID := int64(2)

	currentUser := &models.User{
		ID:   currentUserID,
		Role: "admin",
	}

	targetUser := &models.User{
		ID:   targetUserID,
		Role: "admin",
	}

	mockRepo.On("GetByID", ctx, currentUserID).Return(currentUser, nil)
	mockRepo.On("GetByID", ctx, targetUserID).Return(targetUser, nil)

	err := adminService.MakeAdmin(ctx, targetUserID, currentUserID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "already admin")
	mockRepo.AssertNotCalled(t, "Update")
	mockRepo.AssertExpectations(t)
}

func TestAdminService_RemoveAdmin_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	adminService := NewAdminService(mockRepo)

	ctx := context.Background()
	currentUserID := int64(1)
	targetUserID := int64(2)

	currentUser := &models.User{
		ID:   currentUserID,
		Role: "admin",
	}

	targetUser := &models.User{
		ID:   targetUserID,
		Role: "admin",
	}

	mockRepo.On("GetByID", ctx, currentUserID).Return(currentUser, nil)
	mockRepo.On("GetByID", ctx, targetUserID).Return(targetUser, nil)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(u *models.User) bool {
		return u.ID == targetUserID && u.Role == "user"
	})).Return(nil)

	err := adminService.RemoveAdmin(ctx, targetUserID, currentUserID)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAdminService_RemoveAdmin_TargetNotAdmin(t *testing.T) {
	mockRepo := new(MockUserRepository)
	adminService := NewAdminService(mockRepo)

	ctx := context.Background()
	currentUserID := int64(1)
	targetUserID := int64(2)

	currentUser := &models.User{
		ID:   currentUserID,
		Role: "admin",
	}

	targetUser := &models.User{
		ID:   targetUserID,
		Role: "user",
	}

	mockRepo.On("GetByID", ctx, currentUserID).Return(currentUser, nil)
	mockRepo.On("GetByID", ctx, targetUserID).Return(targetUser, nil)

	err := adminService.RemoveAdmin(ctx, targetUserID, currentUserID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not an admin")
	mockRepo.AssertNotCalled(t, "Update")
	mockRepo.AssertExpectations(t)
}

func TestAdminService_GetUsers_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	adminService := NewAdminService(mockRepo)

	ctx := context.Background()
	page := 1
	pageSize := 20

	passwordHash := "hashed_password"
	email1 := "user1@example.com"
	returnedEmail1 := &email1
	email2 := "user2@example.com"
	returnedEmail2 := &email2
	testUsers := []models.User{
		{
			ID:           1,
			Email:        returnedEmail1,
			PasswordHash: &passwordHash,
			Role:         "user",
		},
		{
			ID:           2,
			Email:        returnedEmail2,
			PasswordHash: &passwordHash,
			Role:         "admin",
		},
	}
	total := int64(2)

	mockRepo.On("GetUsers", ctx, page, pageSize).Return(testUsers, total, nil)

	response, err := adminService.GetUsers(ctx, page, pageSize)

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, total, response.Total)
	assert.Equal(t, page, response.Page)
	assert.Equal(t, pageSize, response.PageSize)
	assert.Equal(t, 1, response.TotalPages)
	assert.Len(t, response.Data, 2)

	for _, user := range response.Data {
		assert.Nil(t, user.PasswordHash)
	}

	mockRepo.AssertExpectations(t)
}

func TestAdminService_GetUsers_RepositoryError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	adminService := NewAdminService(mockRepo)

	ctx := context.Background()
	page := 1
	pageSize := 20

	mockRepo.On("GetUsers", ctx, page, pageSize).Return(([]models.User)(nil), int64(0), assert.AnError)

	response, err := adminService.GetUsers(ctx, page, pageSize)

	require.Error(t, err)
	assert.Nil(t, response)
	mockRepo.AssertExpectations(t)
}
