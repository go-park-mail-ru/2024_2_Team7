package middleware

import (
	httpErrors "kudago/internal/http/errors"
	"net/http"
	"strings"

	"kudago/internal/http/auth"
	"kudago/internal/http/utils"
	"kudago/internal/models"
)

const (
	SessionToken = "session_token"
)

func AuthWithCSRFMiddleware(whitelist []string, authHandler *auth.AuthHandler, encryptionKey []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie(SessionToken)
		if err == nil {
			session, authenticated := authHandler.CheckSessionMiddleware(r.Context(), cookie.Value)
			if authenticated {
				ctx := utils.SetSessionInContext(r.Context(), session)
				r = r.WithContext(ctx)

				if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
					csrfToken := r.Header.Get("X-CSRF-Token")
					if csrfToken == "" {
						utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrCSRFTokenMissing)
						return
					}

					valid, err := authHandler.CheckCSRFMiddleware(r.Context(), encryptionKey, &session, csrfToken)
					if err != nil || !valid {
						utils.WriteResponse(w, http.StatusForbidden, httpErrors.ErrInvalidCSRFToken)
						return
					}
				} else if r.Method == http.MethodGet {
					session, ok := utils.GetSessionFromContext(r.Context())
					if ok {
						csrfToken, err := authHandler.CreateCSRFMiddleware(r.Context(), encryptionKey, &session)
						if err != nil {
							utils.WriteResponse(w, http.StatusInternalServerError, httpErrors.ErrCSRFTokenGenerationFailed)
							return
						}

						w.Header().Set("X-CSRF-Token", csrfToken)

						ctx := utils.SetCSRFInContext(r.Context(), models.TokenData{CSRFtoken: csrfToken})
						r = r.WithContext(ctx)
					}
				}

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
	})
}
