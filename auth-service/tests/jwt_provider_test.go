package tests

import (
	"testing"
	"time"

	"github.com/natrayanp/GoMicro/auth-service/internal/auth/jwt"

	"github.com/stretchr/testify/assert"
)

func TestJWTProvider_GenerateAndValidate(t *testing.T) {
	provider := jwt.NewProvider("test-secret-key", 15*time.Minute, 7*24*time.Hour)

	// Generate tokens
	tokenPair, err := provider.GenerateTokenPair("user123")
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)

	// Validate access token
	userID, tokenType, err := provider.ValidateToken(tokenPair.AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, "user123", userID)
	assert.Equal(t, "access", tokenType)

	// Validate refresh token
	userID, tokenType, err = provider.ValidateToken(tokenPair.RefreshToken)
	assert.NoError(t, err)
	assert.Equal(t, "user123", userID)
	assert.Equal(t, "refresh", tokenType)
}

func TestJWTProvider_InvalidToken(t *testing.T) {
	provider := jwt.NewProvider("test-secret-key", 15*time.Minute, 7*24*time.Hour)

	// Invalid token
	_, _, err := provider.ValidateToken("invalid.token.here")
	assert.Error(t, err)

	// Wrong secret
	provider2 := jwt.NewProvider("different-secret", 15*time.Minute, 7*24*time.Hour)
	tokenPair, _ := provider.GenerateTokenPair("user123")

	_, _, err = provider2.ValidateToken(tokenPair.AccessToken)
	assert.Error(t, err)
}
