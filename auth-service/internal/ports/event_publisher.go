package ports

import (
	"context"

	"github.com/natrayanp/auth-service/internal/domain"
)

type EventPublisherPort interface {
	PublishUserRegistered(ctx context.Context, user *domain.User) error
	PublishUserLoggedIn(ctx context.Context, userID string) error
	PublishUserLoggedOut(ctx context.Context, userID string) error
	PublishPasswordChanged(ctx context.Context, userID string) error
}
