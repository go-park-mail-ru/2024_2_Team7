package middleware

import (
	"context"
	"net/http"
	"strings"

	"kudago/internal/http/auth"
)

const (
	SessionToken = "session_token"
	SessionKey   = "session"
)

func AuthMiddleware(whitelist []string, authHandler *auth.AuthHandler, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(SessionToken)
		if err == nil {
			session, authenticated := authHandler.Service.CheckSession(r.Context(), cookie.Value)
			if authenticated {
				ctx := context.WithValue(r.Context(), SessionKey, session)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}
		}

		for _, path := range whitelist {
			if strings.HasPrefix(r.URL.Path, path) {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	})
}
