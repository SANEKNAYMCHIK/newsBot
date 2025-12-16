package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DBUrl            string
	ServerPort       string
	JWTSecret        string
	ParserInterval   int
	TelegramBotToken string

	EnableHTTPS    bool
	HTTPSCertFile  string
	HTTPSKeyFile   string
	HTTPSPort      string
	AllowedOrigins []string
}

func Load() *Config {
	return &Config{
		DBUrl:            getEnv("CONN_STR", "localhost"),
		ServerPort:       getEnv("PORT", "8080"),
		JWTSecret:        getEnv("JWT_SECRET", "secret-key"),
		ParserInterval:   getEnvAsInt("RSS_PARSER_INTERVAL", 20),
		TelegramBotToken: getEnv("TOKEN", ""),
		EnableHTTPS:      getEnvAsBool("ENABLE_HTTPS", false),
		HTTPSCertFile:    getEnv("HTTPS_CERT_FILE", "ssl/server.crt"),
		HTTPSKeyFile:     getEnv("HTTPS_KEY_FILE", "ssl/server.key"),
		HTTPSPort:        getEnv("HTTPS_PORT", "8443"),
		AllowedOrigins:   getEnvAsSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}
	return defaultValue
}
