package main

import (
	"github.com/gorilla/mux"
	"kudago/internal/auth"
	"kudago/internal/events"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	// Маршруты
	r.HandleFunc("/register", auth.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", auth.AuthMiddleware(auth.LogoutHandler)).Methods("POST")
	r.HandleFunc("/events", events.GetNewsHandler).Methods("GET")

	handler := auth.CORSMiddleware(r)
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
