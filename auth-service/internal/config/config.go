package config

import (
    "os"
    "strconv"
    "time"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    JWT      JWTConfig
}

type ServerConfig struct {
    Port int
    Host string
}

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    Name     string
}

type JWTConfig struct {
    SecretKey     string
    AccessExpiry  time.Duration
    RefreshExpiry time.Duration
}

func Load() *Config {
    return &Config{
        Server: ServerConfig{
            Port: getEnvAsInt("SERVER_PORT", 50051),
            Host: getEnv("SERVER_HOST", "0.0.0.0"),
        },
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnvAsInt("DB_PORT", 5432),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", "password"),
            Name:     getEnv("DB_NAME", "auth_service"),
        },
        JWT: JWTConfig{
            SecretKey:     getEnv("JWT_SECRET", "default-secret-key"),
            AccessExpiry:  getEnvAsDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
            RefreshExpiry: getEnvAsDuration("JWT_REFRESH_EXPIRY", 168*time.Hour), // 7 days
        },
    }
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value, exists := os.LookupEnv(key); exists {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
    if value, exists := os.LookupEnv(key); exists {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}