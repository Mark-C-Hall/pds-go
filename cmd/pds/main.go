package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mark-c-hall/pds-go/internal/api/handler"
	"github.com/mark-c-hall/pds-go/internal/api/router"
	"github.com/mark-c-hall/pds-go/internal/repository"
	"github.com/mark-c-hall/pds-go/internal/service"
	"github.com/mark-c-hall/pds-go/internal/util"
)

func main() {
	logger := log.New(os.Stdout, "PDS: ", log.LstdFlags)

	db, err := repository.SetupDatabase("./pds.db")
	if err != nil {
		logger.Fatalf("Failed to setup database: %v", err)
	}

	accountRepo := repository.NewSQLAccountRepository(db)
	passwordHasher := util.NewBcryptPasswordHasher()
	accountService := service.NewAccountService(accountRepo, passwordHasher, logger)
	accountHandler := handler.NewAccountHandler(accountService, logger)

	router := router.SetupRouter(accountHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		logger.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Println("Shutdown signal received, initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server gracefully stopped")
}
