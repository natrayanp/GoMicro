package grpc

import (
    "context"
    "fmt"
    "log"
    "net"
    "strconv"
    
    "auth-service/internal/config"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

// GrpcServer manages the gRPC server lifecycle
type GrpcServer struct {
    config     *config.Config
    grpcServer *grpc.Server
    handler    *GrpcAuthHandler
}

// NewGrpcServer creates a new gRPC server
func NewGrpcServer(cfg *config.Config, handler *GrpcAuthHandler) *GrpcServer {
    // Create gRPC server with interceptors
    server := grpc.NewServer(
        grpc.ChainUnaryInterceptor(
            loggingInterceptor(),
            recoveryInterceptor(),
            metricsInterceptor(),
        ),
    )
    
    return &GrpcServer{
        config:     cfg,
        grpcServer: server,
        handler:    handler,
    }
}

// RegisterService registers the auth service handler
func (s *GrpcServer) RegisterService() {
    pb.RegisterAuthServiceServer(s.grpcServer, s.handler)
}

// Start starts the gRPC server
func (s *GrpcServer) Start() error {
    addr := net.JoinHostPort(s.config.Server.Host, strconv.Itoa(s.config.Server.Port))
    
    lis, err := net.Listen("tcp", addr)
    if err != nil {
        return fmt.Errorf("failed to listen on %s: %w", addr, err)
    }
    
    // Enable reflection for debugging and tools like grpcurl
    reflection.Register(s.grpcServer)
    
    log.Printf("🚀 gRPC server listening on %s", addr)
    
    return s.grpcServer.Serve(lis)
}

// Stop gracefully stops the gRPC server
func (s *GrpcServer) Stop() {
    if s.grpcServer != nil {
        s.grpcServer.GracefulStop()
        log.Println("gRPC server stopped gracefully")
    }
}

// GrpcServer returns the underlying gRPC server instance
func (s *GrpcServer) GrpcServer() *grpc.Server {
    return s.grpcServer
}

// Interceptors
func loggingInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        log.Printf("gRPC method called: %s", info.FullMethod)
        return handler(ctx, req)
    }
}

func recoveryInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
        defer func() {
            if r := recover(); r != nil {
                err = fmt.Errorf("panic recovered: %v", r)
                log.Printf("Panic in gRPC method %s: %v", info.FullMethod, r)
            }
        }()
        
        return handler(ctx, req)
    }
}

func metricsInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        // TODO: Add metrics collection
        return handler(ctx, req)
    }
}