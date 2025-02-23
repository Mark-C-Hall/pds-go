package router

import (
	"net/http"

	"github.com/mark-c-hall/pds-go/internal/api/handler"
	"github.com/mark-c-hall/pds-go/internal/api/middleware"
)

func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// XRPC Routes
	mux.HandleFunc("/xrpc/com.atproto.server.createAccount",
		middleware.RequireJSON(
			middleware.MethodOnly(http.MethodPost,
				handler.HandleCreateAccount)))

	return mux
}
