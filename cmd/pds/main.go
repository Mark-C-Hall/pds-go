// cmd/pds/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mark-c-hall/pds-go/internal/config"
	"github.com/mark-c-hall/pds-go/internal/handlers"
	"github.com/mark-c-hall/pds-go/internal/server"
)

func main() {
	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	// Call run and handle potential error
	if err := run(cfg); err != nil {
		log.Fatalf("Application error: %v\n", err)
	}
}

func run(cfg *config.Config) error {
	// Initialize context that will be used for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create routers
	mainRouter := http.NewServeMux()
	atprotoRouter := http.NewServeMux()

	// Register global routes directly on main router
	server.RegisterRoutes(mainRouter)

	// Initialize handlers with dependencies
	accountHandler := handlers.NewAccountHandler()

	// Get account routes
	accountRouter := accountHandler.RegisterRoutes()

	// Mount account routes in the com.atproto.server namespace
	server.MountRouter(atprotoRouter, "/xrpc/com.atproto.server", accountRouter)

	// Mount the atproto router on the main router
	mainRouter.Handle("/", atprotoRouter)

	// Create and start the server
	srv := server.NewServer(cfg, mainRouter)

	// Create error channels for monitoring server lifecycle
	serverErrors := make(chan error, 1)
	shutdownComplete := make(chan struct{})

	// Run server in a go routine
	go func() {
		log.Printf("Server starting on %s:%s\n", cfg.Server.Host, cfg.Server.Port)
		serverErrors <- srv.ListenAndServe()
	}()

	// Start a goroutine to handle graceful shutdown
	go func() {
		// Wait for context cancellation (triggered by signal)
		<-ctx.Done()
		log.Println("Shutdown signal received")

		// Create a deadline for graceful shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			if err := srv.Close(); err != nil {
				log.Printf("Forced shutdown failed: %v", err)
			}
		}

		log.Println("Server shutdown complete")
		close(shutdownComplete)
	}()

	// Wait for either server error or complete shutdown
	select {
	case err := <-serverErrors:
		if err != http.ErrServerClosed {
			return err
		}
	case <-shutdownComplete:
		// Shutdown completed successfully
	}

	return nil
}
