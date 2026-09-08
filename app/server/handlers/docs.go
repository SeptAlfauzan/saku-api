package handlers

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Saku API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.32.14/swagger-ui.css">
  <style>
    body { margin: 0; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5.32.14/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5.32.14/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function () {
      window.ui = SwaggerUIBundle({
        url: "./openapi.yaml",
        dom_id: "#swagger-ui",
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        layout: "StandaloneLayout"
      });
    };
  </script>
</body>
</html>
`

func DocsIndex(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(swaggerUIHTML)
}

func OpenAPISpec(c *fiber.Ctx) error {
	data, err := os.ReadFile("openapi.yaml")
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "failed to load openapi spec")
	}
	c.Set("Content-Type", "application/yaml; charset=utf-8")
	return c.Send(data)
}
