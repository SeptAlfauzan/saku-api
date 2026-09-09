# =========================
# Build stage
# =========================
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

# Download dependencies first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build for the target device
RUN CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build -o /app/server .


# =========================
# Runtime stage
# =========================
FROM --platform=$TARGETPLATFORM alpine:3.22

WORKDIR /app

# CA certificates are useful for HTTPS requests to Gemini
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .
COPY --from=builder /app/openapi.yaml .

# Your Fiber port
EXPOSE 3000

CMD ["./server"]
