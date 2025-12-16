package domain

import "time"

type RefreshToken struct {
    ID        string     `json:"id"`
    UserID    string     `json:"user_id"`
    TokenHash string     `json:"-"`
    ExpiresAt time.Time  `json:"expires_at"`
    CreatedAt time.Time  `json:"created_at"`
    RevokedAt *time.Time `json:"revoked_at,omitempty"`
}

func (rt *RefreshToken) IsExpired() bool {
    return time.Now().After(rt.ExpiresAt)
}

func (rt *RefreshToken) IsRevoked() bool {
    return rt.RevokedAt != nil
}

func (rt *RefreshToken) IsValid() bool {
    return !rt.IsExpired() && !rt.IsRevoked()
}