package server

import (
	"net/http"
	"github.com/gorilla/mux"
	"mini/internal/config"
)

// NewRouter creates a new HTTP router with the provided handlers and configuration.
func NewRouter(h *Handlers, cfg *config.App) http.Handler {
	r := mux.NewRouter()
	r.Use(authMiddleware(cfg))

	r.HandleFunc("/", h.RedirectRoot).Methods("GET")
	r.HandleFunc("/view", h.ViewDir).Methods("GET")
	r.HandleFunc("/get", h.GetFile).Methods("GET")
	r.HandleFunc("/upload", h.Upload).Methods("POST")
	r.HandleFunc("/delete", h.Delete).Methods("POST")

	return r
}

// authMiddleware : HTTP Basic-Auth using bcrypt hashed secret in cfg
func authMiddleware(cfg *config.App) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok {
				unauthorized(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="mini"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

