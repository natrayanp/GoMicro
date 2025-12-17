package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/natrayanp/auth-service/internal/domain"
)

type Provider struct {
	secretKey     string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewJWTProvider(secret string, accessExp, refreshExp time.Duration) *Provider {
	return &Provider{
		secretKey:     secret,
		accessExpiry:  accessExp,
		refreshExpiry: refreshExp,
	}
}

func (p *Provider) GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(p.accessExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(p.secretKey))
}

func (p *Provider) GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(p.refreshExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(p.secretKey))
}

func (p *Provider) GenerateTokenPair(userID string) (*domain.TokenPair, error) {
	access, err := p.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refresh, err := p.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int64(p.accessExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (p *Provider) ValidateToken(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(p.secretKey), nil
	})

	if err != nil || !token.Valid {
		return "", "", domain.ErrInvalidToken
	}

	claims := token.Claims.(jwt.MapClaims)
	sub := claims["sub"].(string)

	return sub, tokenString, nil
}

func (p *Provider) HashRefreshToken(token string) string {
	return token // replace with real hashing later
}
