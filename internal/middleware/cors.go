package middleware

import (
	"net/http"
)

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
			"127.0.0.1:8080",
			"http://127.0.0.1",
			"http://127.0.0.1:8080",
		}

		// Проверка на разрешенные домены
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
				break
			}
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
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
