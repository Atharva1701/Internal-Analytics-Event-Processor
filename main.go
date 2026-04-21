package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Structured logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load environment variables
	cfg := LoadConfig()

	// Connect to PostgreSQL
	dbPool, err := NewDBPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer dbPool.Close()

	server := &Server{DB: dbPool}

	// Register HTTP handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/ingest", server.IngestEvent)

	// Health endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Analytics: count total events
	mux.HandleFunc("/analytics/count", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var count int
		err := dbPool.QueryRow(ctx, `SELECT COUNT(*) FROM analytics_events`).Scan(&count)
		if err != nil {
			http.Error(w, "failed to fetch count", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"total_events": %d}`, count)))
	})

	// HTTP server with timeouts
	httpServer := &http.Server{
		Addr:              cfg.ServerAddr,
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("server listening on %s", cfg.ServerAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = httpServer.Shutdown(ctx)
}
