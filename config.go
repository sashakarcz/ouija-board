package main

import (
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	ServerAddr     string
	OllamaURL      string
	OllamaModel    string
	OllamaTimeout  time.Duration
	MaxHistorySize int
	MaxTokens      int
	RateLimit      int
	EnableOTEL     bool
	OTELEndpoint   string
}

// LoadConfig loads configuration from environment variables with sensible defaults
func LoadConfig() *Config {
	return &Config{
		ServerAddr:     getEnv("SERVER_ADDR", "0.0.0.0:8080"),
		OllamaURL:      getEnv("OLLAMA_URL", "http://localhost:11434/api/generate"),
		OllamaModel:    getEnv("OLLAMA_MODEL", "qwen3"),
		OllamaTimeout:  getDurationEnv("OLLAMA_TIMEOUT", 30*time.Second),
		MaxHistorySize: getIntEnv("MAX_HISTORY_SIZE", 1000),
		MaxTokens:      getIntEnv("MAX_TOKENS", 10),
		RateLimit:      getIntEnv("RATE_LIMIT", 10), // requests per second
		EnableOTEL:     getBoolEnv("ENABLE_OTEL", false),
		OTELEndpoint:   getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317"),
	}
}

// Helper functions to read environment variables with defaults
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
