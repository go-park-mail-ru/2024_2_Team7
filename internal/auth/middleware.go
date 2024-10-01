package auth

import (
	"log"
	"net/http"
	"strings"
)

func (h *Handler) AuthMiddleware(whitelist []string, authHandler *Handler, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, path := range whitelist {
			if strings.HasPrefix(r.URL.Path, path) {
				next.ServeHTTP(w, r)
				return
			}
		}

		_, authenticated := authHandler.SessionDb.CheckSession(r)
		if !authenticated {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем CORS-заголовки
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{
			"http://localhost",       // Клиент на порту 80
			"http://vyhodnoy.online", // Другой разрешенный домен
			"http://37.139.40.252",
			"http://37.139.40.252:8080",
			"37.139.40.252",
		}

		// Проверка на разрешенные домены
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
				break
			}
		}
		w.Header().Set("Access-Control-Allow-Origin", "139")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Обработка preflight-запросов
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Запрос:", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
