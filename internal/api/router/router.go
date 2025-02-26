package router

import (
	"net/http"

	"github.com/mark-c-hall/pds-go/internal/api/handler"
	"github.com/mark-c-hall/pds-go/internal/api/middleware"
)

func SetupRouter(accountHandler *handler.AccountHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// XRPC Routes
	mux.HandleFunc("/xrpc/com.atproto.server.createAccount",
		middleware.RequireJSON(
			middleware.MethodOnly(http.MethodPost,
				accountHandler.HandleCreateAccount)))

	return mux
}
