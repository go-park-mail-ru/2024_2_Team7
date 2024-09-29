package main

import (
	"log"
	"net/http"

	"kudago/config"
	"kudago/internal/auth"
	"kudago/internal/events"
	"kudago/internal/users"
	"kudago/session"

	"github.com/gorilla/mux"
)

func main() {
	port := config.LoadConfig()

	r := mux.NewRouter()
	userDB := users.NewUserDB()
	sessionDB := session.NewSessionDB()
	eventsDB := events.NewEventDB()

	authHandler := &auth.AuthHandler{
		UserDB:    *userDB,
		SessionDb: *sessionDB,
	}

	eventHandler := &events.EventHandler{
		EventDB: *eventsDB,
	}

	whitelist := []string{
		"/login",
		"/register",
		"/events",
	}

	r.HandleFunc("/register", authHandler.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", authHandler.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", authHandler.LogoutHandler).Methods("POST")
	r.HandleFunc("/session", authHandler.CheckSessionHandler).Methods("GET")
	r.HandleFunc("/events", eventHandler.GetAllEventsHandler).Methods("GET")
	r.HandleFunc("/events/{tag}", eventHandler.GetEventsByTagHandler).Methods("GET")


	handlerWithCORS := auth.CORSMiddleware(r)

	handler := authHandler.AuthMiddleware(whitelist, *authHandler, handlerWithCORS)

	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
