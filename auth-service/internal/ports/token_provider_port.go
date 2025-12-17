package ports

import (
	"github.com/natrayanp/auth-service/internal/domain"
)

type TokenProviderPort interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	GenerateTokenPair(userID string) (*domain.TokenPair, error)
	ValidateToken(tokenString string) (string, string, error)
	HashRefreshToken(token string) string
}
