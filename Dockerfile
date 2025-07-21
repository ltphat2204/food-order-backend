# syntax=docker/dockerfile:1

# --- Build stage ---
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build REST API binary
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/rest ./cmd/api/main.go
# Build WebSocket server binary
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/socket ./cmd/socket/main.go

# --- REST API image ---
FROM alpine:3.19 AS rest
WORKDIR /app
COPY --from=builder /app/bin/rest ./rest
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/README.MD ./README.MD
COPY --from=builder /app/docker-compose.dev.yml ./docker-compose.dev.yml
COPY --from=builder /app/internal/shared/config/*.yaml ./config/
EXPOSE 8080
ENTRYPOINT ["/app/rest"]

# --- WebSocket image ---
FROM alpine:3.19 AS socket
WORKDIR /app
COPY --from=builder /app/bin/socket ./socket
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/README.MD ./README.MD
COPY --from=builder /app/docker-compose.dev.yml ./docker-compose.dev.yml
COPY --from=builder /app/internal/shared/config/*.yaml ./config/
EXPOSE 8081
ENTRYPOINT ["/app/socket"] 