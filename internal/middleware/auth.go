package middleware

import (
	"net/http"
	"strings"

	"kudago/internal/http/auth"
	"kudago/internal/repository"
)

func AuthMiddleware(whitelist []string, authHandler *auth.AuthHandler, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, path := range whitelist {
			if strings.HasPrefix(r.URL.Path, path) {
				next.ServeHTTP(w, r)
				return
			}
		}

		cookie, _ := r.Cookie(repository.SessionToken)

		_, authenticated := authHandler.Service.CheckSession(r.Context(), cookie.Value)
		if !authenticated {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
