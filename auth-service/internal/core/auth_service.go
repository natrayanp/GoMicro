package core

import (
    "context"
    "time"
    
    "auth-service/internal/domain"
    "auth-service/internal/ports"
    "auth-service/internal/auth/jwt"
    "golang.org/x/crypto/bcrypt"
)

// AuthService implements the AuthServicePort interface
type AuthService struct {
    userRepo       ports.UserRepository
    tokenRepo      ports.TokenRepository
    tokenProvider  ports.TokenProviderPort
    eventPublisher ports.EventPublisherPort // optional
}

func NewAuthService(
    userRepo ports.UserRepository,
    tokenRepo ports.TokenRepository,
    tokenProvider ports.TokenProviderPort,
    eventPublisher ports.EventPublisherPort,
) *AuthService {
    return &AuthService{
        userRepo:       userRepo,
        tokenRepo:      tokenRepo,
        tokenProvider:  tokenProvider,
        eventPublisher: eventPublisher,
    }
}

// Register implements AuthServicePort.Register
func (s *AuthService) Register(ctx context.Context, email, password string) (*domain.User, error) {
    // Check if user exists
    existing, _ := s.userRepo.GetUserByEmail(ctx, email)
    if existing != nil {
        return nil, domain.ErrUserExists
    }
    
    // Validate input
    if err := domain.ValidateEmail(email); err != nil {
        return nil, err
    }
    
    if err := domain.ValidatePassword(password); err != nil {
        return nil, err
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    // Create user
    user, err := domain.NewUser(email, string(hashedPassword))
    if err != nil {
        return nil, err
    }
    
    // Save to database
    if err := s.userRepo.CreateUser(ctx, user); err != nil {
        return nil, err
    }
    
    // Publish event if publisher exists
    if s.eventPublisher != nil {
        go s.eventPublisher.PublishUserRegistered(context.Background(), user)
    }
    
    return user, nil
}

// Login implements AuthServicePort.Login
func (s *AuthService) Login(ctx context.Context, email, password string) (*jwt.TokenPair, error) {
    // Find user
    user, err := s.userRepo.GetUserByEmail(ctx, email)
    if err != nil {
        return nil, domain.ErrInvalidCredentials
    }
    
    // Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return nil, domain.ErrInvalidCredentials
    }
    
    // Generate tokens
    tokenPair, err := s.tokenProvider.GenerateTokenPair(user.ID)
    if err != nil {
        return nil, err
    }
    
    // Hash and store refresh token
    refreshTokenHash := s.tokenProvider.HashRefreshToken(tokenPair.RefreshToken)
    expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days
    
    if err := s.tokenRepo.CreateRefreshToken(ctx, user.ID, refreshTokenHash, expiresAt); err != nil {
        return nil, err
    }
    
    // Publish event
    if s.eventPublisher != nil {
        go s.eventPublisher.PublishUserLoggedIn(context.Background(), user.ID)
    }
    
    return tokenPair, nil
}

// ValidateToken implements AuthServicePort.ValidateToken
func (s *AuthService) ValidateToken(ctx context.Context, token string) (string, error) {
    userID, tokenType, err := s.tokenProvider.ValidateToken(token)
    if err != nil {
        return "", err
    }
    
    if tokenType != "access" {
        return "", domain.ErrInvalidToken
    }
    
    return userID, nil
}

// RefreshToken implements AuthServicePort.RefreshToken
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*jwt.TokenPair, error) {
    // Validate refresh token
    userID, tokenType, err := s.tokenProvider.ValidateToken(refreshToken)
    if err != nil {
        return nil, err
    }
    
    if tokenType != "refresh" {
        return nil, domain.ErrInvalidToken
    }
    
    // Check if token exists in database and is not revoked
    tokenHash := s.tokenProvider.HashRefreshToken(refreshToken)
    dbToken, err := s.tokenRepo.GetRefreshToken(ctx, tokenHash)
    if err != nil {
        return nil, domain.ErrInvalidToken
    }
    
    if dbToken.IsExpired() {
        return nil, domain.ErrTokenExpired
    }
    
    if dbToken.IsRevoked() {
        return nil, domain.ErrTokenRevoked
    }
    
    // Generate new token pair
    tokenPair, err := s.tokenProvider.GenerateTokenPair(userID)
    if err != nil {
        return nil, err
    }
    
    // Store new refresh token
    newRefreshTokenHash := s.tokenProvider.HashRefreshToken(tokenPair.RefreshToken)
    expiresAt := time.Now().Add(7 * 24 * time.Hour)
    
    if err := s.tokenRepo.CreateRefreshToken(ctx, userID, newRefreshTokenHash, expiresAt); err != nil {
        return nil, err
    }
    
    // Revoke old refresh token
    if err := s.tokenRepo.RevokeRefreshToken(ctx, tokenHash); err != nil {
        // Log error but continue
    }
    
    return tokenPair, nil
}

// RevokeToken implements AuthServicePort.RevokeToken
func (s *AuthService) RevokeToken(ctx context.Context, refreshToken string) error {
    tokenHash := s.tokenProvider.HashRefreshToken(refreshToken)
    
    err := s.tokenRepo.RevokeRefreshToken(ctx, tokenHash)
    if err != nil {
        return err
    }
    
    // Publish event
    if s.eventPublisher != nil {
        // We need to get user ID from token
        userID, _, _ := s.tokenProvider.ValidateToken(refreshToken)
        if userID != "" {
            go s.eventPublisher.PublishUserLoggedOut(context.Background(), userID)
        }
    }
    
    return nil
}

// RevokeAllUserTokens implements AuthServicePort.RevokeAllUserTokens
func (s *AuthService) RevokeAllUserTokens(ctx context.Context, userID string) error {
    return s.tokenRepo.RevokeAllUserTokens(ctx, userID)
}

// GetUserByID implements AuthServicePort.GetUserByID
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
    return s.userRepo.GetUserByID(ctx, userID)
}

// GetUserByEmail implements AuthServicePort.GetUserByEmail
func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
    return s.userRepo.GetUserByEmail(ctx, email)
}

// UpdateUserPassword implements AuthServicePort.UpdateUserPassword
func (s *AuthService) UpdateUserPassword(ctx context.Context, userID, newPassword string) error {
    // Validate password
    if err := domain.ValidatePassword(newPassword); err != nil {
        return err
    }
    
    // Hash new password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    
    // Update in repository
    err = s.userRepo.UpdateUserPassword(ctx, userID, string(hashedPassword))
    if err != nil {
        return err
    }
    
    // Revoke all existing tokens for security
    if err := s.tokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
        // Log but don't fail
    }
    
    // Publish event
    if s.eventPublisher != nil {
        go s.eventPublisher.PublishPasswordChanged(context.Background(), userID)
    }
    
    return nil
}	
// convertJWTPairToPortsPair converts JWTTokenPair to ports.TokenPair
func convertJWTPairToPortsPair(jwtPair *jwt.JWTTokenPair) *ports.TokenPair {
    if jwtPair == nil {
        return nil
    }
    return &ports.TokenPair{
        AccessToken:  jwtPair.AccessToken,
        RefreshToken: jwtPair.RefreshToken,
        ExpiresIn:    jwtPair.ExpiresIn,
        TokenType:    jwtPair.TokenType,
    }
}
