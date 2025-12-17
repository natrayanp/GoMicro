package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthChecker struct {
	db *pgxpool.Pool
}

func NewHealthChecker(db *pgxpool.Pool) *HealthChecker {
	return &HealthChecker{db: db}
}

func (h *HealthChecker) Check(ctx context.Context) error {
	// Check database connection
	if err := h.db.Ping(ctx); err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	return nil
}

func (h *HealthChecker) HTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := h.Check(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"status": "unhealthy",
				"error":  err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy",
		})
	}
}
