package handlers

import (
	"net/http"
	"strconv"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
	NewsService     *services.NewsService
	CategoryService *services.CategoryService
}

func NewNewsHandler(newsService *services.NewsService, categoryService *services.CategoryService) *NewsHandler {
	return &NewsHandler{
		NewsService:     newsService,
		CategoryService: categoryService,
	}
}

func (n *NewsHandler) GetNews(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "user not authorized"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	news, err := n.NewsService.GetNews(c.Request.Context(), userID.(int64), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, news)
}

func (n *NewsHandler) GetNewsByID(c *gin.Context) {
	newsID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid news id"})
		return
	}

	news, err := n.NewsService.GetNewsByID(c.Request.Context(), newsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	if news == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "news not found"})
		return
	}

	c.JSON(http.StatusOK, news)
}

func (n *NewsHandler) GetActiveSources(c *gin.Context) {
	sources, err := n.NewsService.GetActiveSources(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, sources)
}

func (n *NewsHandler) GetCategories(c *gin.Context) {
	categories, err := n.CategoryService.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}
