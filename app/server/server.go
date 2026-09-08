package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/septalfauzan/saku-api/app/server/config"
	"github.com/septalfauzan/saku-api/app/server/datasources"
	"github.com/septalfauzan/saku-api/app/server/handlers"
	"github.com/septalfauzan/saku-api/app/server/infrastuctures/ratelimiter"
	"github.com/septalfauzan/saku-api/app/server/services"
)

func NewServer(datasources *datasources.Datasources, cfg *config.Config) *fiber.App {
	app := fiber.New()
	apiRoutes := app.Group("/api/v1")
	receiptService := services.NewReceiptService(datasources.Remote.GeminiAPI)

	limiter := ratelimiter.NewMemoryLimiter(cfg.RateLimit, time.Duration(cfg.RateLimitWindowMs)*time.Millisecond)
	rlMiddleware := ratelimiter.NewMiddleware(limiter)

	apiRoutes.Post("/ocr", rlMiddleware, handlers.OCRReceiptImage(receiptService))
	app.Get("/docs", handlers.DocsIndex)
	app.Get("/openapi.yaml", handlers.OpenAPISpec)

	return app
}
