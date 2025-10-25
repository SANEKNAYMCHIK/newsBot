package main

import (
	"log"

	"github.com/SANEKNAYMCHIK/newsBot/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// Connecting db

	router := gin.Default()

	// Setting API

	log.Printf("Server starting on http://localhost:%s", cfg.ServerPort)
	router.Run(":" + cfg.ServerPort)
}
