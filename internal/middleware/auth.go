package middleware

import (
	"context"
	"net/http"
	"strings"

	"kudago/internal/http/auth"
	"kudago/internal/models"
)


const (
	SessionToken = "session_token"
	SessionKey = "session"
)

func AuthMiddleware(whitelist []string, authHandler *auth.AuthHandler, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(SessionToken)
		var session *models.Session
		var authenticated bool
		if err == nil {
			session, authenticated = authHandler.Service.CheckSession(r.Context(), cookie.Value)
			if session != nil {
				sessionInfo := models.SessionInfo{
					Session:       *session,
					Authenticated: authenticated,
				}
				ctx := context.WithValue(r.Context(), SessionKey, sessionInfo)
				r = r.WithContext(ctx)
			}
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
