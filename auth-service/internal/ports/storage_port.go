package ports

import (
	"context"
	"time"

	"github.com/natrayanp/GoMicro/auth-service/internal/domain"
)

// UserRepository defines storage operations for users
type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUserPassword(ctx context.Context, userID, newPasswordHash string) error
	DeleteUser(ctx context.Context, id string) error
}

// TokenRepository defines storage operations for tokens
type TokenRepository interface {
	CreateRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RevokeAllUserTokens(ctx context.Context, userID string) error
	GetValidRefreshTokens(ctx context.Context, userID string) ([]*domain.RefreshToken, error)
}
