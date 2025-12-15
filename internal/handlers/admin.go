package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	AdminService    *services.AdminService
	SourceService   *services.SourceService
	CategoryService *services.CategoryService
}

func NewAdminHandler(
	AdminService *services.AdminService,
	SourceService *services.SourceService,
	CategoryService *services.CategoryService,
) *AdminHandler {
	return &AdminHandler{
		AdminService:    AdminService,
		SourceService:   SourceService,
		CategoryService: CategoryService,
	}
}

func (a *AdminHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	users, err := a.AdminService.GetUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (a *AdminHandler) UpdateSource(c *gin.Context) {
	sourceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid source id"})
		return
	}

	var req models.UpdateSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	source, err := a.SourceService.UpdateSource(c.Request.Context(), sourceID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, source)
}

func (a *AdminHandler) DeleteSource(c *gin.Context) {
	sourceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid source id"})
		return
	}

	err = a.SourceService.DeleteSource(c.Request.Context(), sourceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: "source deleted"})
}

func (a *AdminHandler) AddCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	category, err := a.CategoryService.CreateCategory(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (a *AdminHandler) MakeAdmin(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Incorrect user ID"})
		return
	}

	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User unathorized"})
		return
	}

	err = a.AdminService.MakeAdmin(c.Request.Context(), userID, currentUserID.(int64))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{
		Message: fmt.Sprintf("User with ID: %d now is admin", userID),
	})
}

func (a *AdminHandler) RemoveAdmin(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Incorrect user ID"})
		return
	}

	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User unathorized"})
		return
	}

	err = a.AdminService.RemoveAdmin(c.Request.Context(), userID, currentUserID.(int64))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{
		Message: fmt.Sprintf("User with ID: %d now is user", userID),
	})
}
