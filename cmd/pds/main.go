// cmd/pds/main.go
package main

import (
	"log"
	"net/http"

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

	log.Printf("Server starting on %s:%s\n", cfg.Server.Host, cfg.Server.Port)
	log.Fatal(srv.ListenAndServe())
}
