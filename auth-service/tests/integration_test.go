package tests

import (
	"context"
	"testing"

	"github.com/natrayanp/GoMicro/auth-service/internal/config"
	"github.com/natrayanp/GoMicro/auth-service/internal/storage/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestDatabaseIntegration(t *testing.T) {
	ctx := context.Background()

	// Start PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	defer postgresContainer.Terminate(ctx)

	// Get connection details
	host, err := postgresContainer.Host(ctx)
	assert.NoError(t, err)

	port, err := postgresContainer.MappedPort(ctx, "5432")
	assert.NoError(t, err)

	// Test database connection
	cfg := &config.DatabaseConfig{
		Host:     host,
		Port:     port.Int(),
		User:     "test",
		Password: "test",
		Name:     "testdb",
	}

	db, err := postgres.NewConnection(cfg)
	if err != nil {
		t.Skip("Could not connect to test database:", err)
	}
	defer db.Close()

	// Run migrations (simplified)
	_, err = db.Pool.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            email VARCHAR(255) UNIQUE NOT NULL,
            password_hash VARCHAR(255) NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
        )
    `)
	assert.NoError(t, err)

	// Test insert and query
	_, err = db.Pool.Exec(ctx, `
        INSERT INTO users (email, password_hash)
        VALUES ($1, $2)
    `, "test@example.com", "hashed_password")
	assert.NoError(t, err)

	var count int
	err = db.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}
