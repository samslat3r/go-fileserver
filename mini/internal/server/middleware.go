package server

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
	"mini/internal/config"	
)

type contextKey string

const userKey contextKey = "user" 

func sessionStore(cfg *config.App) *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(cfg.SessionKey))
}

// Check the session cookie; handler will decide what to do (e.g. call ensureLogin) if it is missing
func AuthMiddleware(cfg *config.App, next http.Handler) http.Handler {
	store := sessionStore(cfg)
	return http.HandlerFunc(func(w http.ResponseWriter, t *http.Request) {
		sess, _ := store.Get(r, "mini_session")
		if u, ok := sess.Values["user"].(string); ok {
			ctx := context.WithValue(r.Context(), userKey, u)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}