package grpc

import (
	"context"
	"log"

	"github.com/natrayanp/GoMicro/auth-service/internal/core"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	UnimplementedAuthServiceServer
	authService *core.AuthService
}

func NewAuthHandler(authService *core.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("Register request for email: %s", req.Email)

	user, err := h.authService.Register(req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{
		UserId:  user.ID,
		Message: "User registered successfully",
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Printf("Login request for email: %s", req.Email)

	tokenPair, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	}

	return &pb.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

func (h *AuthHandler) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	log.Println("Refresh token request")

	tokenPair, err := h.authService.Refresh(req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid refresh token")
	}

	return &pb.RefreshResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

func (h *AuthHandler) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	log.Println("Validate token request")

	userID, err := h.authService.Validate(req.Token)
	if err != nil {
		return &pb.ValidateResponse{Valid: false}, nil
	}

	return &pb.ValidateResponse{
		Valid:  true,
		UserId: userID,
	}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	log.Println("Logout request")

	err := h.authService.Logout(req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to logout")
	}

	return &pb.LogoutResponse{Success: true}, nil
}
