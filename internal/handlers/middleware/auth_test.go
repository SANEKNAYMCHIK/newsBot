package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	"github.com/SANEKNAYMCHIK/newsBot/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager("test-secret")
	email := "test@example.com"
	// Создаем тестового пользователя и токен
	testUser := &models.User{
		ID:    1,
		Email: &email,
		Role:  "user",
	}

	token, err := jwtManager.GenerateToken(testUser)
	require.NoError(t, err)

	router := gin.Default()
	router.Use(AuthMiddleware(jwtManager))

	router.GET("/protected", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		require.True(t, exists)

		email, exists := c.Get("user_email")
		require.True(t, exists)

		role, exists := c.Get("user_role")
		require.True(t, exists)

		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
			"email":   email,
			"role":    role,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем, что данные пользователя установлены в контекст
	// (уже проверено в хендлере выше)
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager("test-secret")

	router := gin.Default()
	router.Use(AuthMiddleware(jwtManager))
	router.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	// Не устанавливаем заголовок Authorization

	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager("test-secret")

	router := gin.Default()
	router.Use(AuthMiddleware(jwtManager))
	router.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRoleMiddleware_AdminAccess(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	// Middleware, который устанавливает роль "admin"
	router.Use(func(c *gin.Context) {
		c.Set("user_role", "admin")
		c.Next()
	})

	router.Use(RoleMiddleware("admin"))
	router.GET("/admin-only", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleMiddleware_UserNoAccess(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	// Middleware, который устанавливает роль "user"
	router.Use(func(c *gin.Context) {
		c.Set("user_role", "user")
		c.Next()
	})

	router.Use(RoleMiddleware("admin"))
	router.GET("/admin-only", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
}
