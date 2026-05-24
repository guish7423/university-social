package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/guish/university-social/internal/cache"
	"github.com/guish/university-social/internal/config"
	"github.com/guish/university-social/internal/migration"
	"github.com/guish/university-social/internal/router"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}
	if err := migration.Run(db, "migrations"); err != nil {
		log.Printf("Migration warning: %v", err)
	}

	var rdb *cache.Cache
	if cfg.RedisAddr != "" {
		var err error
		rdb, err = cache.New(cfg.RedisAddr)
		if err != nil {
			log.Printf("Redis not available (running without cache): %v", err)
		} else {
			defer rdb.Close()
			log.Printf("Redis connected: %s", cfg.RedisAddr)
		}
	}

	r := router.Setup(db, cfg, rdb)

	srv := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on %s", cfg.ServerAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited gracefully")
}
