package postgres

import (
    "context"
    "time"
    
    "auth-service/internal/domain"
    "auth-service/internal/ports"
    "auth-service/internal/storage/postgres/sqlc"
)

// TokenRepository implements ports.TokenRepository
type TokenRepository struct {
    queries *sqlc.Queries
    db      *DB
}

// NewTokenRepository creates a new PostgreSQL token repository
func NewTokenRepository(db *DB) ports.TokenRepository {
    return &TokenRepository{
        queries: sqlc.New(db.Pool),
        db:      db,
    }
}

// CreateRefreshToken implements ports.TokenRepository.CreateRefreshToken
func (r *TokenRepository) CreateRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
    params := sqlc.CreateRefreshTokenParams{
        UserID:    userID,
        TokenHash: tokenHash,
        ExpiresAt: expiresAt,
    }
    
    _, err := r.queries.CreateRefreshToken(ctx, params)
    return err
}

// GetRefreshToken implements ports.TokenRepository.GetRefreshToken
func (r *TokenRepository) GetRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
    result, err := r.queries.GetRefreshToken(ctx, tokenHash)
    if err != nil {
        return nil, domain.ErrInvalidToken
    }
    
    return &domain.RefreshToken{
        ID:        result.ID,
        UserID:    result.UserID,
        TokenHash: result.TokenHash,
        ExpiresAt: result.ExpiresAt,
        CreatedAt: result.CreatedAt,
        RevokedAt: result.RevokedAt,
    }, nil
}

// RevokeRefreshToken implements ports.TokenRepository.RevokeRefreshToken
func (r *TokenRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
    return r.queries.RevokeRefreshToken(ctx, tokenHash)
}

// RevokeAllUserTokens implements ports.TokenRepository.RevokeAllUserTokens
func (r *TokenRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
    return r.queries.RevokeAllUserTokens(ctx, userID)
}

// GetValidRefreshTokens implements ports.TokenRepository.GetValidRefreshTokens
func (r *TokenRepository) GetValidRefreshTokens(ctx context.Context, userID string) ([]*domain.RefreshToken, error) {
    results, err := r.queries.GetValidRefreshTokens(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    tokens := make([]*domain.RefreshToken, len(results))
    for i, result := range results {
        tokens[i] = &domain.RefreshToken{
            ID:        result.ID,
            UserID:    result.UserID,
            TokenHash: result.TokenHash,
            ExpiresAt: result.ExpiresAt,
            CreatedAt: result.CreatedAt,
            RevokedAt: result.RevokedAt,
        }
    }
    
    return tokens, nil
}