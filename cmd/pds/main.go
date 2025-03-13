package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mark-c-hall/pds-go/internal/api/handler"
	"github.com/mark-c-hall/pds-go/internal/api/router"
	"github.com/mark-c-hall/pds-go/internal/config"
	"github.com/mark-c-hall/pds-go/internal/repository"
	"github.com/mark-c-hall/pds-go/internal/service"
	"github.com/mark-c-hall/pds-go/internal/util"
)

func main() {
	logger := log.New(os.Stdout, "PDS: ", log.LstdFlags)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Configuration err: %v", err)
	}

	// Setup database
	db, err := repository.SetupDatabase(cfg)
	if err != nil {
		logger.Fatalf("Failed to setup database: %v", err)
	}

	accountRepo := repository.NewSQLAccountRepository(db)
	passwordHasher := util.NewBcryptPasswordHasher()
	accountService := service.NewAccountService(accountRepo, passwordHasher, logger)
	accountHandler := handler.NewAccountHandler(accountService, logger)

	router := router.SetupRouter(accountHandler)

	server := http.Server{
		Addr:         cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		logger.Printf("Starting server on %s\n", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Println("Shutdown signal received, initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.IdleTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server gracefully stopped")
}
