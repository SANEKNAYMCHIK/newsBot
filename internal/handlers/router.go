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
	refreshService *services.RefreshService,
	jwtManager *auth.JWTManager,
) *gin.Engine {

	router := gin.Default()

	authHandler := NewAuthHandler(authService)
	userHandler := NewUserHandler(userService)
	refreshHandler := NewRefreshHandler(refreshService)
	newsHandler := NewNewsHandler(newsService, sourceService, categoryService)
	subscriptionHandler := NewSubscriptionHandler(subscriptionService)
	adminHandler := NewAdminHandler(adminService, sourceService, categoryService)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/telegram", authHandler.RegisterTelegram)
	}

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		userGroup := protected.Group("/user")
		{
			userGroup.GET("/profile", userHandler.GetProfile)
			userGroup.POST("/refresh", refreshHandler.RequestRefresh)
			userGroup.GET("/refresh/:id", refreshHandler.GetRefreshStatus)
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
			newsGroup.POST("/sources", newsHandler.AddSource)
			newsGroup.GET("/categories", newsHandler.GetCategories)
			newsGroup.GET("/source/:id", newsHandler.GetNewsBySource)
		}

		adminGroup := protected.Group("/admin")
		adminGroup.Use(middleware.RoleMiddleware("admin"))
		{
			adminGroup.POST("/users/:id/make-admin", adminHandler.MakeAdmin)
			adminGroup.POST("/users/:id/remove-admin", adminHandler.RemoveAdmin)

			adminGroup.GET("/users", adminHandler.GetUsers)
			adminGroup.PUT("/sources/:id", adminHandler.UpdateSource)
			adminGroup.DELETE("/sources/:id", adminHandler.DeleteSource)
			adminGroup.POST("/categories", adminHandler.AddCategory)
		}
	}

	return router
}
