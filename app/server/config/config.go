package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GeminiAPIKey    string
	GeminiAPIUrl    string
	GeminiOCRPrompt string
	GeminiOCRModel  string
	Port            string
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}

	config := Config{
		GeminiAPIKey:    os.Getenv("GEMINI_API_KEY"),
		GeminiAPIUrl:    os.Getenv("GEMINI_API_URL"),
		GeminiOCRPrompt: os.Getenv("PROMPT"),
		GeminiOCRModel:  os.Getenv("GEMINI_MODEL"),
		Port:            os.Getenv("PORT"),
	}
	return config, nil
}
