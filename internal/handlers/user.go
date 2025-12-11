package handlers

import (
	"net/http"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (u *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User not unauthorized"})
		return
	}

	user, err := u.UserService.GetProfile(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// func (u *UserHandler) GetSubscriptions(c *gin.Context) {
// 	userID, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User not unauthorized"})
// 		return
// 	}

// 	subscriptions, err := u.UserService.GetSubscriptions(c.Request.Context(), userID.(int64))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, subscriptions)
// }

// func (u *UserHandler) AddSubscription(c *gin.Context) {
// 	userID, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User not unauthorized"})
// 		return
// 	}

// 	var req models.AddSubscriptionRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
// 		return
// 	}

// 	if err := u.UserService.AddSubscription(c.Request.Context(), userID.(int64), req.SourceID); err != nil {
// 		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, models.MessageResponse{Message: "Subscription added successfully"})
// }

// func (u *UserHandler) RemoveSubscription(c *gin.Context) {
// 	userID, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User not unauthorized"})
// 		return
// 	}

// 	sourceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "User not unauthorized"})
// 		return
// 	}

// 	if err := u.UserService.RemoveSubscription(c.Request.Context(), userID.(int64), sourceID); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.MessageResponse{Message: "Subscription removed successfully"})
// }
