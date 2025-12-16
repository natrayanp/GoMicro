package postgres

import (
    "context"
    
    "auth-service/internal/domain"
    "auth-service/internal/ports"
    "auth-service/internal/storage/postgres/sqlc"
)

// UserRepository implements ports.UserRepository
type UserRepository struct {
    queries *sqlc.Queries
    db      *DB
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *DB) ports.UserRepository {
    return &UserRepository{
        queries: sqlc.New(db.Pool),
        db:      db,
    }
}

// CreateUser implements ports.UserRepository.CreateUser
func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
    params := sqlc.CreateUserParams{
        Email:        user.Email,
        PasswordHash: user.PasswordHash,
    }
    
    result, err := r.queries.CreateUser(ctx, params)
    if err != nil {
        return err
    }
    
    user.ID = result.ID
    user.CreatedAt = result.CreatedAt
    user.UpdatedAt = result.UpdatedAt
    
    return nil
}

// GetUserByID implements ports.UserRepository.GetUserByID
func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
    result, err := r.queries.GetUserByID(ctx, id)
    if err != nil {
        return nil, domain.ErrUserNotFound
    }
    
    return &domain.User{
        ID:           result.ID,
        Email:        result.Email,
        PasswordHash: result.PasswordHash,
        CreatedAt:    result.CreatedAt,
        UpdatedAt:    result.UpdatedAt,
    }, nil
}

// GetUserByEmail implements ports.UserRepository.GetUserByEmail
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
    result, err := r.queries.GetUserByEmail(ctx, email)
    if err != nil {
        return nil, domain.ErrUserNotFound
    }
    
    return &domain.User{
        ID:           result.ID,
        Email:        result.Email,
        PasswordHash: result.PasswordHash,
        CreatedAt:    result.CreatedAt,
        UpdatedAt:    result.UpdatedAt,
    }, nil
}

// UpdateUserPassword implements ports.UserRepository.UpdateUserPassword
func (r *UserRepository) UpdateUserPassword(ctx context.Context, userID, newPasswordHash string) error {
    params := sqlc.UpdateUserPasswordParams{
        ID:           userID,
        PasswordHash: newPasswordHash,
    }
    
    return r.queries.UpdateUserPassword(ctx, params)
}

// DeleteUser implements ports.UserRepository.DeleteUser
func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
    return r.queries.DeleteUser(ctx, id)
}