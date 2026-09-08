package remote

import (
	"net/http"
)

type RemoteDatasource struct {
	GeminiAPI GeminiAPI
}

func NewRemoteDatasource(apiKey, apiURL, ocrPrompt, ocrModel string) *RemoteDatasource {
	httpClient := NewHTTPClient(http.DefaultClient, map[string]string{
		"x-goog-api-key": apiKey,
		"Content-Type":   "application/json",
		"Api-Revision":   "2026-05-20",
	})

	return &RemoteDatasource{
		GeminiAPI: NewGeminiAPI(
			httpClient,
			apiURL,
			ocrPrompt,
			ocrModel,
		),
	}
}
