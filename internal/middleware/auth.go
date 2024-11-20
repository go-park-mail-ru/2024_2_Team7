package middleware

import (
	"context"
	"net/http"
	"strings"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/utils"
	"kudago/internal/models"
)

const (
	SessionToken = "session_token"
	SessionKey   = "session"
)

var whitelist = []string{
	"/login",
	"/register",
	"/events",
	"/static",
	"/session",
	"/logout",
	"/docs",
	"/categories",
	"/swagger",
	"/profile",
}

type sessionChecker interface {
	CheckSession(ctx context.Context, cookie string) (models.Session, error)
}

func AuthMiddleware(sessionChecker sessionChecker, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(SessionToken)
		if err == nil {
			session, err := sessionChecker.CheckSession(r.Context(), cookie.Value)
			if err == nil {
				ctx := utils.SetSessionInContext(r.Context(), session)
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

		utils.WriteResponse(w, http.StatusUnauthorized, httpErrors.ErrUnauthorized)
		return
	})
}
