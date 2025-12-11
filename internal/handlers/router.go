package handlers

import (
	"github.com/SANEKNAYMCHIK/newsBot/internal/handlers/middleware"
	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
	"github.com/SANEKNAYMCHIK/newsBot/pkg/auth"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	authService *services.AuthService,
	userService *services.UserService,
	newsService *services.NewsService,
	categoryService *services.CategoryService,
	subscriptionService *services.SubscriptionService,
	sourceService *services.SourceService,
	adminService *services.AdminService,
	jwtManager *auth.JWTManager,
) *gin.Engine {

	router := gin.Default()

	authHandler := NewAuthHandler(authService)
	userHandler := NewUserHandler(userService)
	newsHandler := NewNewsHandler(newsService, sourceService, categoryService)
	subscriptionHandler := NewSubscriptionHandler(subscriptionService)
	adminHandler := NewAdminHandler(adminService, sourceService, categoryService)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		userGroup := protected.Group("/user")
		{
			userGroup.GET("/profile", userHandler.GetProfile)
			// userGroup.PUT("/profile", userHandler.UpdateProfile)
		}

		subscriptionGroup := protected.Group("/user/subscriptions")
		{
			subscriptionGroup.GET("/", subscriptionHandler.GetSubscriptions)
			subscriptionGroup.POST("/", subscriptionHandler.AddSubscription)
			subscriptionGroup.DELETE("/:id", subscriptionHandler.RemoveSubscription)
		}

		newsGroup := protected.Group("/news")
		{
			newsGroup.GET("/", newsHandler.GetNews)
			newsGroup.GET("/:id", newsHandler.GetNewsByID)
			newsGroup.GET("/sources", newsHandler.GetActiveSources)
			newsGroup.GET("/categories", newsHandler.GetCategories)
		}

		adminGroup := protected.Group("/admin")
		adminGroup.Use(middleware.RoleMiddleware("admin"))
		{
			adminGroup.GET("/users", adminHandler.GetUsers)
			adminGroup.POST("/sources", adminHandler.AddSource)
			adminGroup.PUT("/sources/:id", adminHandler.UpdateSource)
			adminGroup.DELETE("/sources/:id", adminHandler.DeleteSource)
			adminGroup.POST("/categories", adminHandler.AddCategory)
		}
	}

	return router
}
