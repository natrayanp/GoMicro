package domain

import "errors"

var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrUserExists         = errors.New("user already exists")
    ErrUserNotFound       = errors.New("user not found")
    ErrInvalidToken       = errors.New("invalid token")
    ErrTokenExpired       = errors.New("token expired")
    ErrInvalidEmail       = errors.New("invalid email address")
    ErrPasswordTooShort   = errors.New("password must be at least 8 characters")
    ErrTokenRevoked       = errors.New("token has been revoked")
)