package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/natrayanp/auth-service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	config     *config.Config
	grpcServer *grpc.Server
}

func NewServer(cfg *config.Config) *Server {
	// Create gRPC server with interceptors
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingInterceptor(),
			recoveryInterceptor(),
		),
	)

	return &Server{
		config:     cfg,
		grpcServer: server,
	}
}

func (s *Server) Start() error {
	addr := net.JoinHostPort(s.config.Server.Host, strconv.Itoa(s.config.Server.Port))

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	log.Printf("Starting gRPC server on %s", addr)

	// Enable reflection for debugging
	reflection.Register(s.grpcServer)

	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}

func loggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Printf("gRPC method: %s", info.FullMethod)
		return handler(ctx, req)
	}
}

func recoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic recovered: %v", r)
				log.Printf("Panic in %s: %v", info.FullMethod, r)
			}
		}()

		return handler(ctx, req)
	}
}
