package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/natrayanp/GoMicro/auth-service/internal/domain"
	"github.com/natrayanp/GoMicro/auth-service/internal/ports"
	"github.com/natrayanp/GoMicro/auth-service/internal/storage/postgres/sqlc"
)

type UserRepository struct {
	queries *sqlc.Queries
	db      *DB
}

func NewUserRepository(db *DB) ports.UserRepository {
	return &UserRepository{
		queries: sqlc.New(db.Pool),
		db:      db,
	}
}

// ------------------------------
// CREATE USER
// ------------------------------

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	params := sqlc.CreateUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}

	result, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	// Convert pgtype → domain
	user.ID = result.ID.String()
	user.CreatedAt = result.CreatedAt.Time
	user.UpdatedAt = result.UpdatedAt.Time

	return nil
}

// ------------------------------
// GET USER BY ID
// ------------------------------

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	// Convert string → pgtype.UUID
	uid := pgtype.UUID{}
	_ = uid.Scan(id)

	result, err := r.queries.GetUserByID(ctx, uid)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return &domain.User{
		ID:           result.ID.String(),
		Email:        result.Email,
		PasswordHash: result.PasswordHash,
		CreatedAt:    result.CreatedAt.Time,
		UpdatedAt:    result.UpdatedAt.Time,
	}, nil
}

// ------------------------------
// GET USER BY EMAIL
// ------------------------------

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	result, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return &domain.User{
		ID:           result.ID.String(),
		Email:        result.Email,
		PasswordHash: result.PasswordHash,
		CreatedAt:    result.CreatedAt.Time,
		UpdatedAt:    result.UpdatedAt.Time,
	}, nil
}

// ------------------------------
// UPDATE PASSWORD
// ------------------------------

func (r *UserRepository) UpdateUserPassword(ctx context.Context, userID, newPasswordHash string) error {
	uid := pgtype.UUID{}
	_ = uid.Scan(userID)

	params := sqlc.UpdateUserPasswordParams{
		ID:           uid,
		PasswordHash: newPasswordHash,
	}

	return r.queries.UpdateUserPassword(ctx, params)
}

// ------------------------------
// DELETE USER
// ------------------------------

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	uid := pgtype.UUID{}
	_ = uid.Scan(id)

	return r.queries.DeleteUser(ctx, uid)
}
