package router

import (
	"net/http"

	"github.com/mark-c-hall/pds-go/internal/api/handler"
)

func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// XRPC Routes
	mux.HandleFunc("/xrpc/com.atproto.server.createAccount", handler.HandleCreateAccount)

	return mux
}
