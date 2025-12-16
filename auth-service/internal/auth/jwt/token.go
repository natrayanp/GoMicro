package jwt

import "time"

// Renamed to JWTTokenPair to avoid conflict
type JWTTokenPair struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    ExpiresIn    int64     `json:"expires_in"`
    TokenType    string    `json:"token_type"`
}

func NewTokenPair(accessToken, refreshToken string, expiresIn int64) *JWTTokenPair {
    return &JWTTokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    expiresIn,
        TokenType:    "Bearer",
    }
}

type RefreshToken struct {
    ID        string
    UserID    string
    TokenHash string
    ExpiresAt time.Time
    CreatedAt time.Time
    RevokedAt *time.Time
}

func (rt *RefreshToken) IsExpired() bool {
    return time.Now().After(rt.ExpiresAt)
}

func (rt *RefreshToken) IsRevoked() bool {
    return rt.RevokedAt != nil
}
