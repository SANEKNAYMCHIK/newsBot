package main

import (
	"context"
	"log"

	"github.com/SANEKNAYMCHIK/newsBot/internal/config"
	"github.com/SANEKNAYMCHIK/newsBot/internal/database"
	"github.com/SANEKNAYMCHIK/newsBot/internal/repositories"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()
	db, err := database.NewPostgres(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %w", err)
	}
	defer db.Close()
	log.Println("Database connection established successfully")

	userRepo := repositories.NewUserRepository(db.Pool)

	router := gin.Default()

	// Setting API

	log.Printf("Server starting on http://localhost:%s", cfg.ServerPort)
	router.Run(":" + cfg.ServerPort)
}
