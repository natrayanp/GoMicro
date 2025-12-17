package ports

import (
	"context"

	"github.com/natrayanp/auth-service/internal/domain"
)

type AuthServicePort interface {
	Register(ctx context.Context, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.TokenPair, error)
	ValidateToken(ctx context.Context, token string) (string, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error)
	RevokeToken(ctx context.Context, refreshToken string) error
	RevokeAllUserTokens(ctx context.Context, userID string) error
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUserPassword(ctx context.Context, userID, newPassword string) error
}
