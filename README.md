# Auth Service with SQLC

A production-ready authentication service with SQLC from day one.

## Features
- ✅ JWT-based authentication
- ✅ Refresh token rotation
- ✅ SQLC for type-safe SQL
- ✅ gRPC API
- ✅ PostgreSQL database
- ✅ Health checks
- ✅ Docker support

## Quick Start

## 🐳 Docker Commands

### Quick Start (Production)
```bash
# Build and start everything
docker-compose up -d

# Check logs
docker-compose logs -f auth-service

# Stop everything
docker-compose down


### Development (Hot Reload)
# Start with hot reload
docker compose -f docker-compose.dev.yml up -d

# View logs
docker compose -f docker-compose.dev.yml logs -f

# Run tests
docker-compose -f docker-compose.test.yml up


###Utility Commands

# Access database
make db-shell

# Rebuild services
docker-compose build --no-cache

# Clean everything
make clean



##Health Checks
Service: http://localhost:8080/health

gRPC: localhost:50051


## 🎯 **10. Complete Working Script `start.sh`**

#!/bin/bash
echo "🚀 Auth Service - Docker Only Setup"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker Desktop."
    exit 1
fi

# Check for .env file
if [ ! -f .env ]; then
    echo "📄 Creating .env file from example..."
    cp .env.example .env
fi

echo "🐳 Building and starting services..."
docker-compose up -d --build

echo "⏳ Waiting for services to be ready..."
sleep 10

# Check if services are healthy
if curl -s http://localhost:8080/health | grep -q "healthy"; then
    echo "✅ Services are up and running!"
    echo "🌐 Health check: http://localhost:8080/health"
    echo "🔌 gRPC endpoint: localhost:50051"
    echo "📊 PostgreSQL: localhost:5432"
else
    echo "⚠️  Services may still be starting..."
    echo "📋 Check logs: docker-compose logs -f auth-service"
fi



🔧 How It Works Now:
Single Command to Start Everything:
bash
# Just run this ONE command:
docker-compose up -d

# Or with the script:
chmod +x start.sh
./start.sh