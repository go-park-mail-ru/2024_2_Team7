package main

import (
	"github.com/gorilla/mux"
	"log"
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

	// Запуск сервера
	log.Fatal(http.ListenAndServe(":8080", r))
}
