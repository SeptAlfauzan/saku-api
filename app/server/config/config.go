package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	GeminiAPIKey      string
	GeminiAPIUrl      string
	GeminiOCRPrompt   string
	GeminiOCRModel    string
	Port              string
	RateLimit         int
	RateLimitWindowMs int
}

func intEnv(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func LoadConfig() (Config, error) {
	// Missing .env isn't fatal: fall back to OS env vars and defaults below.
	_ = godotenv.Load()

	return Config{
		GeminiAPIKey:      os.Getenv("GEMINI_API_KEY"),
		GeminiAPIUrl:      os.Getenv("GEMINI_API_URL"),
		GeminiOCRPrompt:   os.Getenv("PROMPT"),
		GeminiOCRModel:    os.Getenv("GEMINI_MODEL"),
		Port:              os.Getenv("PORT"),
		RateLimit:         intEnv("RATE_LIMIT", 5),
		RateLimitWindowMs: intEnv("RATE_LIMIT_WINDOW_MS", 60000),
	}, nil
}
