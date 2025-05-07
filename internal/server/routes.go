// internal/server/routes.go
package server

import (
	"net/http"
	"strings"
)

// MountRouter mounts a sub-router onto the main router with a namespace
// Keeps the dotted notation format for endpoints
func MountRouter(main *http.ServeMux, namespace string, sub *http.ServeMux) {
	// For each pattern in the sub router, we need to register it with the namespace prefix
	// This is a bit tricky since Go's http.ServeMux doesn't expose registered patterns
	// We'll use a proxy handler to intercept requests

	// Ensure namespace ends with a dot if not empty
	if namespace != "" && !strings.HasSuffix(namespace, ".") {
		namespace = namespace + "."
	}

	// Handle requests for this namespace
	main.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Extract the path but remove leading slash
		path := strings.ToLower(r.URL.Path)

		// Check if the path starts with our namespace
		if strings.HasPrefix(path, namespace) {
			// Remove the namespace from the path to get the endpoint
			endpoint := strings.TrimPrefix(path, namespace)

			// Update the request path to just the endpoint
			r.URL.Path = "/" + strings.ToLower(endpoint)

			// Let the sub router handle it
			sub.ServeHTTP(w, r)
			return
		}

		// Not handled by this router, continue to next handler
		http.NotFound(w, r)
	})
}

// RegisterRoutes registers global routes
func RegisterRoutes(router *http.ServeMux) {
	// Register global AT Protocol routes here
	router.HandleFunc("/com.atproto.health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}
