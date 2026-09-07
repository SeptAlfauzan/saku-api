package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// /RAW GEMINIT RESPONSE
type InteractionResponse struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	Usage       Usage  `json:"usage"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
	ServiceTier string `json:"service_tier"`
	Steps       []Step `json:"steps"`
	Object      string `json:"object"`
	Model       string `json:"model"`
}

type Usage struct {
	TotalTokens                int               `json:"total_tokens"`
	TotalInputTokens           int               `json:"total_input_tokens"`
	InputTokensByModality      []TokenByModality `json:"input_tokens_by_modality"`
	TotalCachedTokens          int               `json:"total_cached_tokens"`
	TotalOutputTokens          int               `json:"total_output_tokens"`
	TotalToolUseTokens         int               `json:"total_tool_use_tokens"`
	TotalThoughtTokens         int               `json:"total_thought_tokens"`
	RawPromptToken             int               `json:"raw_prompt_token"`
	ModelInvocationTokenCounts []ModelInvocation `json:"model_invocation_token_counts"`
}

type TokenByModality struct {
	Modality string `json:"modality"`
	Tokens   int    `json:"tokens"`
}

type ModelInvocation struct {
	PromptTokensDetails     []TokenByModality `json:"prompt_tokens_details"`
	CandidatesTokensDetails []TokenByModality `json:"candidates_tokens_details"`
}

type Step struct {
	Signature string    `json:"signature,omitempty"`
	Type      string    `json:"type"`
	Content   []Content `json:"content,omitempty"`
}

type Content struct {
	Text string `json:"text,omitempty"`
	Type string `json:"type"`
}

// FORMATTED RESPONSE
type ReceiptResponse struct {
	MerchantName    *string       `json:"merchant_name"`
	TransactionDate *string       `json:"transaction_date"`
	TransactionTime *string       `json:"transaction_time"`
	Currency        *string       `json:"currency"`
	Subtotal        *float64      `json:"subtotal"`
	Tax             *float64      `json:"tax"`
	Discount        *float64      `json:"discount"`
	Total           *float64      `json:"total"`
	PaymentMethod   *string       `json:"payment_method"`
	Items           []ReceiptItem `json:"items"`
}

type ReceiptItem struct {
	Name       string  `json:"name"`
	Quantity   float64 `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found", err)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, from SAKU!")
	})

	///proxy API to gemini API
	app.Post("/api/v1/proxy/ocr", func(c *fiber.Ctx) error {
		var input struct {
			ImageBase64 string `json:"image"`
			MimeType    string `json:"mime_type"`
		}

		if err := c.BodyParser(&input); err != nil {
			log.Println("error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		var apiURL string = os.Getenv("GEMINI_API_URL")
		var apiKey string = os.Getenv("GEMINI_API_KEY")
		var prompt string = os.Getenv("PROMPT")

		payload := map[string]interface{}{
			"model": "gemini-3.5-flash-lite",
			"input": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": prompt,
				},
				map[string]interface{}{
					"type":      "image",
					"data":      input.ImageBase64,
					"mime_type": input.MimeType,
				},
			},
		}

		jsonBody, err := json.Marshal(payload)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err,
			})
		}

		req, err := http.NewRequest(
			http.MethodPost,
			apiURL, bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Println("error", err)
			return c.Status(500).JSON(fiber.Map{
				"error": "failed to create request",
			})
		}

		req.Header.Add("x-goog-api-key", apiKey)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Api-Revision", "2026-05-20")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err,
			})
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			return c.Status(resp.StatusCode).Send(body)
		}

		var response InteractionResponse
		errDecode := json.NewDecoder(resp.Body).Decode(&response)
		if errDecode != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": errDecode,
			})
		}

		rawFormattedJSON := response.Steps[1].Content[0].Text
		var receipt ReceiptResponse
		errDecodeFormattedJSON := json.Unmarshal([]byte(rawFormattedJSON), &receipt)
		if errDecodeFormattedJSON != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": errDecodeFormattedJSON,
			})
		}

		if len(receipt.Items) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Image don't have items on it! Mostlikely not a receipt image",
				"image": input.ImageBase64,
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"data": receipt,
		})
	})

	app.Listen(":" + port)
}
