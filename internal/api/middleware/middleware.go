package middleware

import (
	"net/http"
	"slices"

	"github.com/mark-c-hall/pds-go/internal/api/httputil"
)

func RequireJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !slices.Contains(r.Header["Content-Type"], "application/json") {
			httputil.RespondWithError(w, "InvalidRequest",
				"Content-Type must be application/json",
				http.StatusBadRequest)
			return
		}
		next(w, r)
	}
}

func MethodOnly(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			httputil.RespondWithError(w, "InvalidRequest",
				"Method Not Supported",
				http.StatusBadRequest)
			return
		}
		next(w, r)
	}
}
