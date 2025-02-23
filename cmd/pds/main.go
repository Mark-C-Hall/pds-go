package main

import (
	"log"
	"net/http"

	"github.com/mark-c-hall/pds-go/internal/api/router"
)

func main() {
	router := router.SetupRouter()

	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Fatal(srv.ListenAndServe())
}
