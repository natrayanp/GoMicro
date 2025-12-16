package domain

import (
    "regexp"
    "time"
    "unicode/utf8"
)

type User struct {
    ID           string    `json:"id"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func NewUser(email, password string) (*User, error) {
    if err := ValidateEmail(email); err != nil {
        return nil, err
    }
    
    if err := ValidatePassword(password); err != nil {
        return nil, err
    }
    
    return &User{
        Email:     email,
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
    }, nil
}

func ValidateEmail(email string) error {
    if email == "" {
        return ErrInvalidEmail
    }
    
    if !emailRegex.MatchString(email) {
        return ErrInvalidEmail
    }
    
    if utf8.RuneCountInString(email) > 255 {
        return ErrInvalidEmail
    }
    
    return nil
}

func ValidatePassword(password string) error {
    if utf8.RuneCountInString(password) < 8 {
        return ErrPasswordTooShort
    }
    
    return nil
}