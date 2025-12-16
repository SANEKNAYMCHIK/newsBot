package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SANEKNAYMCHIK/newsBot/internal/config"
	"github.com/SANEKNAYMCHIK/newsBot/internal/database"
	"github.com/SANEKNAYMCHIK/newsBot/internal/handlers"
	"github.com/SANEKNAYMCHIK/newsBot/internal/repositories"
	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
	"github.com/SANEKNAYMCHIK/newsBot/internal/worker"
	"github.com/SANEKNAYMCHIK/newsBot/pkg/auth"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()
	db, err := database.NewPostgres(ctx, cfg.DBUrl)
	if err != nil {
		log.Printf("Failed to connect database: %s", err)
	}
	defer db.Close()
	log.Println("Database connection established successfully")

	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	userRepo := repositories.NewUserRepository(db.Pool)
	newsRepo := repositories.NewNewsRepository(db.Pool)
	sourceRepo := repositories.NewSourceRepository(db.Pool)
	categoryRepo := repositories.NewCategoryRepository(db.Pool)
	subscriptionRepo := repositories.NewSubscriptionRepository(db.Pool)

	rssParser := services.NewRssParser(10)
	rssService := services.NewRssService(sourceRepo, newsRepo, rssParser)

	refreshService := services.NewRefreshService(
		rssService,
		subscriptionRepo,
		5,
		100,
		3*time.Minute,
	)
	go refreshService.Start(context.Background())
	newsWorker := worker.NewNewsWorker(rssService, time.Duration(cfg.ParserInterval)*time.Minute)

	go func() {
		log.Println("Starting RSS news worker...")
		newsWorker.Start(context.Background())
	}()

	defer func() {
		newsWorker.Stop()
		log.Println("Workers stopped")
	}()

	authService := services.NewAuthService(userRepo, jwtManager)
	userService := services.NewUserService(userRepo)
	newsService := services.NewNewsService(newsRepo, sourceRepo, subscriptionRepo)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, sourceRepo)
	sourceService := services.NewSourceService(sourceRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	adminService := services.NewAdminService(userRepo)

	router := handlers.NewRouter(authService, userService, newsService, categoryService, subscriptionService, sourceService, adminService, refreshService, jwtManager)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.ServerPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")

}
