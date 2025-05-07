package server

import (
	"fmt"
	"net/http"

	"github.com/mark-c-hall/pds-go/internal/config"
)

func NewServer(cfg *config.Config, handler *http.ServeMux) *http.Server {
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}
