package handlers

import (
	"net/http"
	"strconv"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	SubscriptionService *services.SubscriptionService
}

func NewSubscriptionHandler(SubscriptionService *services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{SubscriptionService: SubscriptionService}
}

func (s *SubscriptionHandler) GetSubscriptions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	subscriptions, err := s.SubscriptionService.GetUserSubscriptions(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

func (s *SubscriptionHandler) AddSubscription(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req models.AddSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	err := s.SubscriptionService.AddSubscription(c.Request.Context(), userID.(int64), req.SourceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.MessageResponse{Message: "subscription added"})
}

func (s *SubscriptionHandler) RemoveSubscription(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	sourceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid source id"})
		return
	}

	err = s.SubscriptionService.RemoveSubscription(c.Request.Context(), userID.(int64), sourceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: "subscription removed"})
}
