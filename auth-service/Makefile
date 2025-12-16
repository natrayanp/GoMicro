.PHONY: up down build dev logs test clean

# Production
up:
	docker-compose up -d --build

down:
	docker-compose down -v

logs:
	docker-compose logs -f auth-service

build:
	docker-compose build

# Development
dev:
	docker-compose -f docker-compose.dev.yml up -d --build

dev-down:
	docker-compose -f docker-compose.dev.yml down -v

dev-logs:
	docker-compose -f docker-compose.dev.yml logs -f auth-service

# Testing
test:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit

# Database
db-shell:
	docker-compose exec postgres psql -U postgres -d auth_service

# Cleanup
clean:
	rm -rf bin/ tmp/ *.pb.go internal/storage/postgres/sqlc/
	docker-compose down -v --rmi all
	docker system prune -f