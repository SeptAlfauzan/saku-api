package remote

import (
	"net/http"

	"github.com/septalfauzan/saku-api/app/server/config"
)

type RemoteDatasource struct {
	GeminiAPI GeminiAPI
}

func NewRemoteDatasource(cfg *config.Config) *RemoteDatasource {
	httpClient := NewHTTPClient(http.DefaultClient, map[string]string{
		"x-goog-api-key": cfg.GeminiAPIKey,
		"Content-Type":   "application/json",
		"Api-Revision":   "2026-05-20",
	})

	return &RemoteDatasource{
		GeminiAPI: NewGeminiAPI(
			httpClient,
			cfg.GeminiAPIUrl,
			cfg.GeminiOCRPrompt,
			cfg.GeminiOCRModel,
		),
	}
}
