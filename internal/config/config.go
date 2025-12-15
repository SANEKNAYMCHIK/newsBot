package config

import "os"

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Redis name;port;...

	TelegramBotToken string
	ServerPort       string
}

func Load() *Config {
	return &Config{
		DBHost: getEnv("DB_HOST", "localhost"),
		// ...
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
