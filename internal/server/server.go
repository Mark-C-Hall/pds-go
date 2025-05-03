package server

import (
	"net/http"

	"github.com/mark-c-hall/pds-go/internal/handlers"
)

func NewServer() *http.Server {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", handlers.CreateAccountHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return srv
}
