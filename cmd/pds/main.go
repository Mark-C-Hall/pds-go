package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mark-c-hall/pds-go/internal/api/handler"
	"github.com/mark-c-hall/pds-go/internal/api/router"
	"github.com/mark-c-hall/pds-go/internal/service"
)

func main() {
	logger := log.New(os.Stdout, "PDS: ", log.LstdFlags)

	accountService := service.NewAccountService()

	accountHandler := handler.NewAccountHandler(accountService, logger)

	mux := router.SetupRouter(accountHandler)

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	logger.Println("Starting server on :8080")
	log.Fatal(srv.ListenAndServe())
}
