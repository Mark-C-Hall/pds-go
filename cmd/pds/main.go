package main

import (
	"log"

	"github.com/mark-c-hall/pds-go/internal/server"
)

func main() {
	srv := server.NewServer()
	log.Fatal(srv.ListenAndServe())
}
