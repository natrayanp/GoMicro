# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache \
    protoc \
    protobuf-dev \
    make \
    git \
    bash

# Copy go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Install code generation tools
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate protobuf code
RUN protoc --go_out=. --go-grpc_out=. proto/*.proto

# Generate SQLC code
RUN sqlc generate

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/auth-service ./cmd

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    curl

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/auth-service /app/
COPY --from=builder /app/.env.example /app/.env

# Create non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 50051 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["/app/auth-service"]