package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/septalfauzan/saku-api/app/server/datasources"
	"github.com/septalfauzan/saku-api/app/server/handlers"
	"github.com/septalfauzan/saku-api/app/server/services"
)

func NewServer(datasources *datasources.Datasources) *fiber.App {
	app := fiber.New()
	apiRoutes := app.Group("/api/v1")
	receiptService := services.NewReceiptService(datasources.Remote.GeminiAPI)
	apiRoutes.Post("/ocr", handlers.OCRReceiptImage(receiptService))

	app.Get("/docs", handlers.DocsIndex)
	app.Get("/openapi.yaml", handlers.OpenAPISpec)

	return app
}
