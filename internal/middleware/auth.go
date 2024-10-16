package middleware

import (
	"context"
	"net/http"
	"strings"

	"kudago/internal/http/auth"
	"kudago/internal/models"
)

func AuthMiddleware(whitelist []string, authHandler *auth.AuthHandler, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var session *models.Session
		var authenticated bool
		cookie, err := r.Cookie(auth.SessionToken)

		if err == nil {
			session, authenticated = authHandler.Service.CheckSession(r.Context(), cookie.Value)
			ctx := context.WithValue(r.Context(), auth.SessionKey, session)
			r = r.WithContext(ctx)
		}

		for _, path := range whitelist {
			if strings.HasPrefix(r.URL.Path, path) {
				next.ServeHTTP(w, r)
				return
			}
		}

		if !authenticated {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
