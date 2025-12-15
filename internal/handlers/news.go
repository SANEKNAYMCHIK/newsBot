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
	SourceService   *services.SourceService
	CategoryService *services.CategoryService
}

func NewNewsHandler(
	newsService *services.NewsService,
	sourceService *services.SourceService,
	categoryService *services.CategoryService,
) *NewsHandler {
	return &NewsHandler{
		NewsService:     newsService,
		SourceService:   sourceService,
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
	sources, err := n.SourceService.GetActiveSources(c.Request.Context())
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

func (n *NewsHandler) AddSource(c *gin.Context) {
	var req models.CreateSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	source, err := n.SourceService.CreateSource(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, source)
}

func (n *NewsHandler) GetNewsBySource(c *gin.Context) {
	sourceIDStr := c.Param("id")
	sourceID, err := strconv.ParseInt(sourceIDStr, 10, 64)
	if err != nil || sourceID <= 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Incorrect source ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	userID, _ := c.Get("user_id")
	response, err := n.NewsService.GetNewsBySource(c.Request.Context(), sourceID, userID.(int64), page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
