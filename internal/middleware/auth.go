package middleware

import (
	"fmt"
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
		// Проверяем сессионный токен
		cookie, err := r.Cookie(SessionToken)
		if err == nil {
			session, authenticated := authHandler.CheckSessionMiddleware(r.Context(), cookie.Value)
			if authenticated {
				ctx := utils.SetSessionInContext(r.Context(), session)
				r = r.WithContext(ctx)

				w.Header().Set("X-Test-Header", "test-value")
				if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
					csrfToken := r.Header.Get("X-CSRF-Token")
					fmt.Println("CSRF Token got from header:", csrfToken)
					if csrfToken == "" {
						utils.WriteResponse(w, http.StatusForbidden, map[string]string{"error": "CSRF token missing"})
						return
					}

					valid, err := authHandler.CheckCSRFMiddleware(r.Context(), encryptionKey, &session, csrfToken)
					if err != nil || !valid {
						utils.WriteResponse(w, http.StatusForbidden, map[string]string{"error": "Invalid CSRF token"})
						return
					}
				} else if r.Method == http.MethodGet {
					session, ok := utils.GetSessionFromContext(r.Context())
					if ok {
						csrfToken, err := authHandler.CreateCSRFMiddleware(r.Context(), encryptionKey, &session)
						if err != nil {
							utils.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate CSRF token"})
							return
						}
						fmt.Println("Generated CSRF Token:", csrfToken)

						w.Header().Set("X-CSRF-Token", csrfToken)
						fmt.Println("Setting X-CSRF-Token in header:", w.Header().Get("X-CSRF-Token"))

						// Далее проверьте, что токен установлен
						for k, v := range w.Header() {
							fmt.Printf("Header after setting CSRF: %s: %s\n", k, v)
						}

						ctx := utils.SetCSRFInContext(r.Context(), models.TokenData{CSRFtoken: csrfToken})
						r = r.WithContext(ctx)
					}
				}

				// Переходим к следующему обработчику
				next.ServeHTTP(w, r)
				return
			}
		}

		// Если маршрут в white list, пропускаем проверку
		for _, path := range whitelist {
			if strings.HasPrefix(r.URL.Path, path) {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Отправляем ошибку, если запрос не авторизован
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
