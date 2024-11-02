package middleware

import (
	"kudago/internal/models"
	"net/http"

	"kudago/internal/http/auth"
	"kudago/internal/http/utils"
)

func CSRFMiddleware(next http.Handler, authHandler *auth.AuthHandler, encryptionKey []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
			csrfToken := r.Header.Get("X-CSRF-Token")
			if csrfToken == "" {
				utils.WriteResponse(w, http.StatusForbidden, map[string]string{"error": "CSRF token missing"})
				return
			}

			// Проверка сессии и токена CSRF
			session, ok := utils.GetSessionFromContext(r.Context())
			if !ok {
				utils.WriteResponse(w, http.StatusUnauthorized, map[string]string{"error": "Session not found"})
				return
			}

			// Проверка токена CSRF
			valid, err := authHandler.CheckCSRFMiddleware(r.Context(), encryptionKey, &session, csrfToken)
			if err != nil || !valid {
				utils.WriteResponse(w, http.StatusForbidden, map[string]string{"error": "Invalid CSRF token"})
				return
			}
		}

		// Генерация нового CSRF токена для GET запросов
		if r.Method == http.MethodGet {
			session, ok := utils.GetSessionFromContext(r.Context())
			if ok {
				csrfToken, err := authHandler.CreateCSRFMiddleware(r.Context(), encryptionKey, &session)
				if err != nil {
					utils.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate CSRF token"})
					return
				}
				// Устанавливаем CSRF токен в заголовке ответа
				w.Header().Set("X-CSRF-Token", csrfToken)
				// Добавляем токен в контекст запроса
				ctx := utils.SetCSRFInContext(r.Context(), models.TokenData{CSRFtoken: csrfToken})
				r = r.WithContext(ctx)
			}
		}

		// Переходим к следующему обработчику
		next.ServeHTTP(w, r)
	})
}
