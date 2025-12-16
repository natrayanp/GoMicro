package jwt

import (
    "crypto/sha256"
    "encoding/hex"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "auth-service/internal/domain"
)

type Provider struct {
    secretKey     []byte
    accessExpiry  time.Duration
    refreshExpiry time.Duration
}

func NewProvider(secretKey string, accessExpiry, refreshExpiry time.Duration) *Provider {
    return &Provider{
        secretKey:     []byte(secretKey),
        accessExpiry:  accessExpiry,
        refreshExpiry: refreshExpiry,
    }
}

func (p *Provider) GenerateAccessToken(userID string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(p.accessExpiry).Unix(),
        "iat":     time.Now().Unix(),
        "type":    "access",
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(p.secretKey)
}

func (p *Provider) GenerateRefreshToken(userID string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(p.refreshExpiry).Unix(),
        "iat":     time.Now().Unix(),
        "type":    "refresh",
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(p.secretKey)
}

func (p *Provider) GenerateTokenPair(userID string) (*JWTTokenPair, error) {
    accessToken, err := p.GenerateAccessToken(userID)
    if err != nil {
        return nil, err
    }
    
    refreshToken, err := p.GenerateRefreshToken(userID)
    if err != nil {
        return nil, err
    }
    
    return NewTokenPair(accessToken, refreshToken, int64(p.accessExpiry.Seconds())), nil
}

func (p *Provider) ValidateToken(tokenString string) (string, string, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return p.secretKey, nil
    })
    
    if err != nil {
        return "", "", domain.ErrInvalidToken
    }
    
    if !token.Valid {
        return "", "", domain.ErrInvalidToken
    }
    
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return "", "", domain.ErrInvalidToken
    }
    
    exp, ok := claims["exp"].(float64)
    if !ok || time.Now().Unix() > int64(exp) {
        return "", "", domain.ErrTokenExpired
    }
    
    userID, ok := claims["user_id"].(string)
    if !ok {
        return "", "", domain.ErrInvalidToken
    }
    
    tokenType, ok := claims["type"].(string)
    if !ok {
        return "", "", domain.ErrInvalidToken
    }
    
    return userID, tokenType, nil
}

func (p *Provider) HashRefreshToken(token string) string {
    hash := sha256.Sum256([]byte(token))
    return hex.EncodeToString(hash[:])
}
