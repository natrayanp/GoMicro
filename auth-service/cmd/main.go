package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/natrayanp/GoMicro/auth-service/internal/adapters/grpc"
	"github.com/natrayanp/GoMicro/auth-service/internal/api/health"
	"github.com/natrayanp/GoMicro/auth-service/internal/auth/jwt"
	"github.com/natrayanp/GoMicro/auth-service/internal/config"
	"github.com/natrayanp/GoMicro/auth-service/internal/core"
	"github.com/natrayanp/GoMicro/auth-service/internal/storage/postgres"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Println("Configuration loaded")

	// Setup database connection
	db, err := postgres.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connected")

	// Setup repositories (implement ports)
	userRepo := postgres.NewUserRepository(db)
	tokenRepo := postgres.NewTokenRepository(db)

	// Setup JWT provider (implements TokenProviderPort)
	jwtProvider := jwt.NewJWTProvider(
		cfg.JWT.SecretKey,
		cfg.JWT.AccessExpiry,
		cfg.JWT.RefreshExpiry,
	)

	// Setup auth service (implements AuthServicePort)
	// Note: eventPublisher is nil for now, can be added later
	authService := core.NewAuthService(userRepo, tokenRepo, jwtProvider, nil)
	log.Println("Core services initialized")

	// Setup gRPC handler (adapter)
	grpcHandler := grpc.NewGrpcAuthHandler(authService)

	// Setup gRPC server
	grpcServer := grpc.NewGrpcServer(cfg, grpcHandler)
	grpcServer.RegisterService()

	// Setup health checks
	healthChecker := health.NewHealthChecker(db.Pool)

	// Start health check HTTP server
	go startHealthServer(healthChecker)

	// Start gRPC server
	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	waitForShutdown(grpcServer)

	log.Println("Server shutdown complete")
}

func startHealthServer(healthChecker *health.HealthChecker) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthChecker.HTTPHandler())
	mux.HandleFunc("/ready", healthChecker.HTTPHandler())

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Starting health server on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Health server failed: %v", err)
	}
}

func waitForShutdown(grpcServer *grpc.GrpcServer) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Println("Shutting down server...")

	grpcServer.Stop()

	log.Println("Server stopped gracefully")
}
