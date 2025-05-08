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

	// Create a channel to listen for server errors
	serverErrors := make(chan error, 1)

	// Run server in a go routine
	go func() {
		log.Printf("Server starting on %s:%s\n", cfg.Server.Host, cfg.Server.Port)
		serverErrors <- srv.ListenAndServe()
	}()

	// Create a channel to listen for signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until an os signal or error is reached
	select {
	case err := <-serverErrors:
		log.Fatalf("Error with server: %v\n", err)

	case sig := <-shutdown:
		log.Printf("Shutdown signal received: %v\n", sig)

		// Create a deadline for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			if err := srv.Close(); err != nil {
				log.Printf("Forced shutdown failed: %v", err)
			}
		}

		log.Println("Server shutdown complete")
	}
}
