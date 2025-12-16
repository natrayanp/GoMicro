package tests

import (
    "testing"
    "time"
    
    "auth-service/internal/auth/jwt"
    "auth-service/internal/core"
    "auth-service/internal/domain"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock repositories
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
    args := m.Called(email)
    return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id string) (*domain.User, error) {
    args := m.Called(id)
    return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *domain.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserRepository) Delete(id string) error {
    args := m.Called(id)
    return args.Error(0)
}

type MockTokenRepository struct {
    mock.Mock
}

func (m *MockTokenRepository) CreateRefreshToken(userID, tokenHash string, expiresAt time.Time) error {
    args := m.Called(userID, tokenHash, expiresAt)
    return args.Error(0)
}

func (m *MockTokenRepository) GetRefreshToken(tokenHash string) (*domain.RefreshToken, error) {
    args := m.Called(tokenHash)
    return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (m *MockTokenRepository) RevokeRefreshToken(tokenHash string) error {
    args := m.Called(tokenHash)
    return args.Error(0)
}

func (m *MockTokenRepository) RevokeAllUserTokens(userID string) error {
    args := m.Called(userID)
    return args.Error(0)
}

func TestAuthService_Register(t *testing.T) {
    mockUserRepo := new(MockUserRepository)
    mockTokenRepo := new(MockTokenRepository)
    
    jwtProvider := jwt.NewProvider("test-secret", 15*time.Minute, 7*24*time.Hour)
    
    authService := core.NewAuthService(mockUserRepo, mockTokenRepo, jwtProvider)
    
    // Test registration
    mockUserRepo.On("FindByEmail", "test@example.com").Return((*domain.User)(nil), domain.ErrUserNotFound)
    mockUserRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)
    
    user, err := authService.Register("test@example.com", "password123")
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test@example.com", user.Email)
    
    mockUserRepo.AssertExpectations(t)
}