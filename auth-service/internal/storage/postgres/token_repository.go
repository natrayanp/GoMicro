package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/natrayanp/auth-service/internal/domain"
	"github.com/natrayanp/auth-service/internal/ports"
	"github.com/natrayanp/auth-service/internal/storage/postgres/sqlc"
)

type TokenRepository struct {
	queries *sqlc.Queries
	db      *DB
}

func NewTokenRepository(db *DB) ports.TokenRepository {
	return &TokenRepository{
		queries: sqlc.New(db.Pool),
		db:      db,
	}
}

// ------------------------------
// CREATE
// ------------------------------

func (r *TokenRepository) CreateRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	// Convert userID → pgtype.UUID
	uid := pgtype.UUID{}
	_ = uid.Scan(userID)

	// Convert expiresAt → pgtype.Timestamptz
	exp := pgtype.Timestamptz{
		Time:  expiresAt,
		Valid: true,
	}

	params := sqlc.CreateRefreshTokenParams{
		UserID:    uid,
		TokenHash: tokenHash,
		ExpiresAt: exp,
	}

	_, err := r.queries.CreateRefreshToken(ctx, params)
	return err
}

// ------------------------------
// GET SINGLE
// ------------------------------

func (r *TokenRepository) GetRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	result, err := r.queries.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	// Convert pgtype → domain types
	id := result.ID.String()
	userID := result.UserID.String()

	expiresAt := result.ExpiresAt.Time
	createdAt := result.CreatedAt.Time

	var revokedAt *time.Time
	if result.RevokedAt.Valid {
		revokedAt = &result.RevokedAt.Time
	}

	return &domain.RefreshToken{
		ID:        id,
		UserID:    userID,
		TokenHash: result.TokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
		RevokedAt: revokedAt,
	}, nil
}

// ------------------------------
// REVOKE SINGLE
// ------------------------------

func (r *TokenRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	return r.queries.RevokeRefreshToken(ctx, tokenHash)
}

// ------------------------------
// REVOKE ALL USER TOKENS
// ------------------------------

func (r *TokenRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	uid := pgtype.UUID{}
	_ = uid.Scan(userID)

	return r.queries.RevokeAllUserTokens(ctx, uid)
}

// ------------------------------
// GET VALID TOKENS FOR USER
// ------------------------------

func (r *TokenRepository) GetValidRefreshTokens(ctx context.Context, userID string) ([]*domain.RefreshToken, error) {
	uid := pgtype.UUID{}
	_ = uid.Scan(userID)

	results, err := r.queries.GetValidRefreshTokens(ctx, uid)
	if err != nil {
		return nil, err
	}

	tokens := make([]*domain.RefreshToken, len(results))
	for i, result := range results {
		id := result.ID.String()
		u := result.UserID.String()

		expiresAt := result.ExpiresAt.Time
		createdAt := result.CreatedAt.Time

		var revokedAt *time.Time
		if result.RevokedAt.Valid {
			revokedAt = &result.RevokedAt.Time
		}

		tokens[i] = &domain.RefreshToken{
			ID:        id,
			UserID:    u,
			TokenHash: result.TokenHash,
			ExpiresAt: expiresAt,
			CreatedAt: createdAt,
			RevokedAt: revokedAt,
		}
	}

	return tokens, nil
}
