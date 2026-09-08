# Swagger/OpenAPI Docs for saku-api

Date: 2026-09-08

## Goal

Add self-hosted interactive API documentation (Swagger UI) for the existing saku-api service. The API currently exposes a single endpoint, `POST /api/v1/ocr`, and has no OpenAPI spec or docs tooling. This change is additive: it documents the existing endpoint without modifying handlers, domain types, or request/response behavior.

## Requirements

- Serve an interactive Swagger UI docs page at `/docs`.
- Serve an OpenAPI 3.0 spec file at `/openapi.yaml`.
- The docs must fully document the existing `POST /api/v1/ocr` endpoint, including request schema, success response, and error responses.
- Docs reflect the actual Go structs used by the API.
- No handler logic, domain types, or routing of the API itself changes. Only additive docs routes are introduced.
- Swagger UI renderer assets (HTML/CSS/JS) load from the unpkg CDN; the OpenAPI spec is served locally by the Fiber app.

## Decisions

Recorded during brainstorming:

- **Scoped to "API docs UI only"** — no code generation, no `oapi-codegen`, no swaggo.
- **UI renderer:** Swagger UI.
- **Delivery:** self-hosted spec + docs route on the Fiber app.
- **Assets:** Swagger UI HTML/CSS/JS loaded from unpkg CDN at runtime.
- **Route placement:** root-level `/docs` and `/openapi.yaml` (not version-namespaced).

## Architecture

```
+-------------+        GET /docs        +----------------------+
|   Browser   | ----------------------> | Fiber app: HTML page |
+-------------+                         |  (loads swagger-ui    |
      |                                 |   from unpkg CDN)    |
      | GET /openapi.yaml               +----------------------+
      +---------------------------------> Fiber app: static yaml |
                                          +----------------------+
      |                                        
      | POST /api/v1/ocr (via "Try it out")
      +---------------------------------> existing OCR handler (unchanged)
```

### Components

1. **`openapi.yaml`** — committed at repo root. OpenAPI 3.0.3 document containing:
   - API title, version, description metadata.
   - A single `paths` entry for `/api/v1/ocr`.
   - Components schemas matching the Go domain types:
     - `OCRRequest` — `image` (string, required), `mime_type` (string).
     - `Receipt` — `merchant_name`, `transaction_date`, `transaction_time`, `currency`, `subtotal`, `tax`, `discount`, `total`, `payment_method`, `items`.
     - `ReceiptItem` — `name`, `quantity`, `unit_price`, `total_price`.
     - `Error` — `error` (string).
   - Success response: `201` with body `{ "data": <Receipt> }`.
   - Error responses: `400` and `500` with body `{ "error": string }`.

2. **`/docs` route** — Fiber GET handler returning an HTML document that:
   - Loads Swagger UI CSS/JS from unpkg.
   - Configures Swagger UI to fetch `./openapi.yaml` (relative).
   - Renders the interactive documentation.

3. **`/openapi.yaml` route** — Fiber GET handler returning the committed `openapi.yaml` with correct `Content-Type` (`application/yaml` or `text/yaml`).

### Routing changes

In `app/server/server.go`, alongside the existing `/api/v1` group:

```go
apiRoutes := app.Group("/api/v1")
receiptService := services.NewReceiptService(datasources.Remote.GeminiAPI)
apiRoutes.Post("/ocr", handlers.OCRReceiptImage(receiptService))

// Additive docs routes
app.Get("/docs", handlers.DocsIndex)      // HTML page
app.Get("/openapi.yaml", handlers.OpenAPISpec) // spec
```

Two new handler functions (likely a new `handlers/docs.go`): `DocsIndex` and `OpenAPISpec`.

## Data Flow

1. User opens `GET /docs`.
2. Fiber returns an HTML page that loads swagger-ui assets from unpkg CDN.
3. Swagger UI fetches `GET /openapi.yaml`, which Fiber serves from the committed spec.
4. Swagger UI renders the endpoint definition.
5. Using "Try it out", the user may `POST /api/v1/ocr`; behavior is identical to the existing API.

## Error Handling

No changes to API error handling. The OpenAPI spec documents the existing responses:
- `400` — invalid request body, missing `image`, or image contains no receipt items (`{ "error": string }`).
- `500` — upstream OCR (Gemini) extraction failure (`{ "error": string }`).

## Testing

- Manual verification: start the server, load `/docs` in a browser, confirm Swagger UI renders and `Try it out` works against `/openapi.yaml`.
- Verify `/openapi.yaml` is served with a YAML content type.
- Verify existing `POST /api/v1/ocr` behavior is unchanged.

## Out of Scope

- Code generation from the spec (`oapi-codegen`).
- Swaggo annotations / `swag init` workflow.
- Authentication/authorization for `/docs`.
- CORS configuration.
- Documentation for any future endpoints beyond OCR.

## Success Criteria

- `GET /docs` returns a working Swagger UI page.
- `GET /openapi.yaml` returns the OpenAPI spec.
- The spec accurately documents `POST /api/v1/ocr` (request, success, and error schemas) matching the Go domain types.
- Existing API behavior is unchanged.
