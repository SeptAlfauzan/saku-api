package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/septalfauzan/saku-api/app/server/domain"
)

type GeminiOCRRequest struct {
	Model string        `json:"model"`
	Input []GeminiInput `json:"input"`
}

type GeminiInput struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	Data     string `json:"data,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
}

type GeminiAPI interface {
	ExtractImageOCR(
		ctx context.Context,
		request domain.OCRRequest,
	) (domain.Receipt, error)
}

type DefaultGeminiAPI struct {
	httpClient *HTTPClient
	apiURL     string
	ocrPrompt  string
	ocrModel   string
}

func NewGeminiAPI(
	httpClient *HTTPClient,
	apiURL string,
	ocrPrompt string,
	ocrModel string,
) GeminiAPI {
	return &DefaultGeminiAPI{
		httpClient: httpClient,
		apiURL:     apiURL,
		ocrPrompt:  ocrPrompt,
		ocrModel:   ocrModel,
	}
}

func (api *DefaultGeminiAPI) ExtractImageOCR(
	ctx context.Context,
	request domain.OCRRequest,
) (domain.Receipt, error) {
	payload := GeminiOCRRequest{
		Model: api.ocrModel,
		Input: []GeminiInput{
			{
				Text: api.ocrPrompt,
				Type: "text",
			},
			{
				Type:     "image",
				Data:     request.ImageBase64,
				MimeType: request.MimeType,
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return domain.Receipt{}, fmt.Errorf(
			"failed to marshal OCR request: %w",
			err,
		)
	}

	res, err := api.httpClient.Do(
		ctx,
		http.MethodPost,
		api.apiURL,
		body,
		nil,
	)
	if err != nil {
		return domain.Receipt{}, fmt.Errorf(
			"failed to execute Gemini request: %w",
			err,
		)
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK ||
		res.StatusCode >= http.StatusMultipleChoices {

		return domain.Receipt{}, fmt.Errorf(
			"Gemini API returned status %d",
			res.StatusCode,
		)
	}

	var geminiResp domain.GeminiResponse

	if err := json.NewDecoder(res.Body).Decode(&geminiResp); err != nil {
		return domain.Receipt{}, fmt.Errorf(
			"failed to decode Gemini response: %w",
			err,
		)
	}

	if len(geminiResp.Steps) < 2 || len(geminiResp.Steps[1].Content) == 0 {
		return domain.Receipt{}, fmt.Errorf(
			"Gemini response missing receipt data in steps",
		)
	}

	rawJSON := geminiResp.Steps[1].Content[0].Text
	var receipt domain.Receipt
	if err := json.Unmarshal([]byte(rawJSON), &receipt); err != nil {
		return domain.Receipt{}, fmt.Errorf(
			"failed to unmarshal receipt from Gemini response: %w",
			err,
		)
	}

	return receipt, nil
}
