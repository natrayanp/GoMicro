package grpc

import (
	"context"
	"log"

	"github.com/natrayanp/auth-service/internal/domain"
	"github.com/natrayanp/auth-service/internal/ports"
	pb "github.com/natrayanp/auth-service/proto/auth/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GrpcAuthHandler adapts gRPC requests to the core service
type GrpcAuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService ports.AuthServicePort
}

// NewGrpcAuthHandler creates a new gRPC handler
func NewGrpcAuthHandler(authService ports.AuthServicePort) *GrpcAuthHandler {
	return &GrpcAuthHandler{
		authService: authService,
	}
}

// Register handles gRPC Register requests
func (h *GrpcAuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("[gRPC] Register request for email: %s", req.Email)

	user, err := h.authService.Register(ctx, req.Email, req.Password)
	if err != nil {
		log.Printf("[gRPC] Register failed: %v", err)
		return nil, mapDomainErrorToGrpc(err)
	}

	return &pb.RegisterResponse{
		UserId:  user.ID,
		Message: "User registered successfully",
	}, nil
}

// Login handles gRPC Login requests
func (h *GrpcAuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Printf("[gRPC] Login request for email: %s", req.Email)

	tokenPair, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		log.Printf("[gRPC] Login failed: %v", err)
		return nil, mapDomainErrorToGrpc(err)
	}

	return &pb.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

// Refresh handles gRPC Refresh requests
func (h *GrpcAuthHandler) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	log.Printf("[gRPC] Refresh token request")

	tokenPair, err := h.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		log.Printf("[gRPC] Refresh failed: %v", err)
		return nil, mapDomainErrorToGrpc(err)
	}

	return &pb.RefreshResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

// Validate handles gRPC Validate requests
func (h *GrpcAuthHandler) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	log.Printf("[gRPC] Validate token request")

	userID, err := h.authService.ValidateToken(ctx, req.Token)
	if err != nil {
		// For validate endpoint, we return valid=false instead of error
		return &pb.ValidateResponse{Valid: false}, nil
	}

	return &pb.ValidateResponse{
		Valid:  true,
		UserId: userID,
	}, nil
}

// Logout handles gRPC Logout requests
func (h *GrpcAuthHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	log.Printf("[gRPC] Logout request")

	if err := h.authService.RevokeToken(ctx, req.RefreshToken); err != nil {
		log.Printf("[gRPC] Logout failed: %v", err)
		return nil, mapDomainErrorToGrpc(err)
	}

	return &pb.LogoutResponse{Success: true}, nil
}

// GetUser handles gRPC GetUser requests
func (h *GrpcAuthHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Printf("[gRPC] GetUser request for ID: %s", req.UserId)

	user, err := h.authService.GetUserByID(ctx, req.UserId)
	if err != nil {
		log.Printf("[gRPC] GetUser failed: %v", err)
		return nil, mapDomainErrorToGrpc(err)
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
	}, nil
}

// mapDomainErrorToGrpc maps domain errors to gRPC status errors
func mapDomainErrorToGrpc(err error) error {
	switch err {
	case domain.ErrInvalidCredentials:
		return status.Error(codes.Unauthenticated, "invalid credentials")
	case domain.ErrUserExists:
		return status.Error(codes.AlreadyExists, "user already exists")
	case domain.ErrUserNotFound:
		return status.Error(codes.NotFound, "user not found")
	case domain.ErrInvalidToken, domain.ErrTokenExpired, domain.ErrTokenRevoked:
		return status.Error(codes.Unauthenticated, "invalid token")
	case domain.ErrInvalidEmail:
		return status.Error(codes.InvalidArgument, "invalid email")
	case domain.ErrPasswordTooShort:
		return status.Error(codes.InvalidArgument, "password too short")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
