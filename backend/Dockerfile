# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build all services
RUN CGO_ENABLED=0 GOOS=linux go build -o fetcher ./cmd/fetcher
RUN CGO_ENABLED=0 GOOS=linux go build -o detector ./cmd/detector
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy all binaries
COPY --from=builder /app/fetcher .
COPY --from=builder /app/detector .
COPY --from=builder /app/api .

# Default command (will be overridden by docker-compose)
CMD ["./api"]