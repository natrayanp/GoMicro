#!/usr/bin/env bash

set -e

SERVICE_NAME=$1

if [ -z "$SERVICE_NAME" ]; then
  echo "Usage: ./bootstrap-service.sh <service-name>"
  exit 1
fi

SERVICE_DIR="./$SERVICE_NAME"

echo "Creating service: $SERVICE_NAME"

# 1. Create folder structure
mkdir -p $SERVICE_DIR/{cmd,internal,proto/$SERVICE_NAME/v1,sql}

# 2. Create main.go
cat > $SERVICE_DIR/cmd/main.go <<EOF
package main

import "fmt"

func main() {
    fmt.Println("$SERVICE_NAME service running...")
}
EOF

# 3. Create .air.toml
cat > $SERVICE_DIR/.air.toml <<EOF
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o tmp/main ./cmd"
bin = "tmp/main"
include_ext = ["go"]
exclude_dir = ["tmp", "vendor"]

[log]
time = true
EOF

# 4. Create Dockerfile.dev
cat > $SERVICE_DIR/Dockerfile.dev <<EOF
FROM golang:1.22-alpine

WORKDIR /app

RUN apk add --no-cache protoc protobuf-dev make git bash curl

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest \\
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \\
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \\
    && go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]
EOF

# 5. Create Dockerfile
cat > $SERVICE_DIR/Dockerfile <<EOF
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o $SERVICE_NAME ./cmd

FROM alpine:3.19
WORKDIR /app

COPY --from=builder /app/$SERVICE_NAME .

CMD ["./$SERVICE_NAME"]
EOF

# 6. Create service-level Makefile
cat > $SERVICE_DIR/Makefile <<EOF
proto:
\tprotoc --go_out=. --go-grpc_out=. proto/$SERVICE_NAME/v1/*.proto

sqlc:
\tsqlc generate -f sql/sqlc.yaml

run:
\tgo run ./cmd

build:
\tgo build -o $SERVICE_NAME ./cmd
EOF

# 7. Insert into docker-compose.dev.yml
if grep -q "services:" docker-compose.dev.yml; then
  echo "Updating docker-compose.dev.yml..."
  sed -i "/services:/a \\
  $SERVICE_NAME:\\
    build:\\
      context: ./$SERVICE_NAME\\
      dockerfile: Dockerfile.dev\\
    working_dir: /app\\
    ports:\\
      - \"8080\"\\
    volumes:\\
      - go-mod-cache:/go/pkg/mod\\
      - go-build-cache:/root/.cache/go-build\\
    networks:\\
      - default" docker-compose.dev.yml
fi

echo "Service $SERVICE_NAME created successfully."
