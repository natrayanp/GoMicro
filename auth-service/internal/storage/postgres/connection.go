package postgres

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "auth-service/internal/config"
    "github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
    Pool *pgxpool.Pool
}

func NewConnection(cfg *config.DatabaseConfig) (*DB, error) {
    connString := fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s?sslmode=disable",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
    )
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    poolConfig, err := pgxpool.ParseConfig(connString)
    if err != nil {
        return nil, fmt.Errorf("failed to parse connection config: %w", err)
    }
    
    pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to create connection pool: %w", err)
    }

    
    // Test connection immediately
    if err := pool.Ping(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    log.Println("Connected to PostgreSQL database")
    return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
    if db.Pool != nil {
        db.Pool.Close()
    }
}