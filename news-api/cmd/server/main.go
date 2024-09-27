package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"news-api/internal/auth"
	"news-api/internal/news"
)

func main() {
	r := mux.NewRouter()

	// Маршруты
	r.HandleFunc("/register", auth.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", auth.AuthMiddleware(auth.LogoutHandler)).Methods("POST")
	r.HandleFunc("/news", news.GetNewsHandler).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:5500"}, // Ваш фронтэнд
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
	})

	// Используйте обработчик CORS
	handler := c.Handler(r)

	// Запуск сервера
	http.ListenAndServe(":8080", handler)
}
