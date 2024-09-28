package main

import (
	"log"
	"net/http"

	"kudago/internal/auth"
	"kudago/internal/events"
	"kudago/internal/users"
	"kudago/session"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	userDB := users.NewUserDB()
	sessionDB := session.NewSessionDB()
	eventsDB := events.NewEventDB()

	authHandler := &auth.Handler{
		UserDB:    *userDB,
		SessionDb: *sessionDB,
		EventDB:   *eventsDB,
	}

	whitelist := []string{
		"/login",
		"/register",
	}

	r.HandleFunc("/register", authHandler.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", authHandler.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", authHandler.LogoutHandler).Methods("POST")
	r.HandleFunc("/", authHandler.GetEventsHandler).Methods("GET")

	handlerWithCORS := auth.CORSMiddleware(r)

	handler := authHandler.AuthMiddleware(whitelist, *authHandler, handlerWithCORS)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
