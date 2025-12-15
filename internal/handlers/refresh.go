package handlers

import (
	"net/http"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
	"github.com/gin-gonic/gin"
)

type RefreshHandler struct {
	refreshService *services.RefreshService
}

func NewRefreshHandler(refreshService *services.RefreshService) *RefreshHandler {
	return &RefreshHandler{refreshService: refreshService}
}

func (r *RefreshHandler) RequestRefresh(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User not unathorized"})
		return
	}
	req, err := r.refreshService.RequestRefresh(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusTooManyRequests, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"request_id": req.ID,
		"status":     req.Status,
		"message":    "Запрос на обновление добавлен в очередь",
	})
}

func (r *RefreshHandler) GetRefreshStatus(c *gin.Context) {
	requestID := c.Param("id")
	if req, ok := r.refreshService.GetRequestStatus(requestID); ok {
		c.JSON(http.StatusOK, req)
		return
	}
	c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "request not found"})
}
