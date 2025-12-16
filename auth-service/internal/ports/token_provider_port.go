package ports

type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int64  `json:"expires_in"`
    TokenType    string `json:"token_type"`
}

type TokenProviderPort interface {
    GenerateAccessToken(userID string) (string, error)
    GenerateRefreshToken(userID string) (string, error)
    GenerateTokenPair(userID string) (*TokenPair, error)
    ValidateToken(tokenString string) (string, string, error)
    HashRefreshToken(token string) string
}
