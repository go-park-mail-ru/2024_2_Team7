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

	authHandler := &auth.Handler{
		UserDB:    *users.NewUserDB(),
		SessionDb: *session.NewSessionDB(),
	}

	eventHandler := &events.Handler{
		EventDB: *events.NewEventDB(),
	}

	whitelist := []string{
		"/login",
		"/register",
		"/events",
		"/static",
		"/session",
	}

	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	r.HandleFunc("/session", authHandler.CheckSession).Methods("GET")
	r.HandleFunc("/events", eventHandler.GetAllEvents).Methods("GET")
	r.HandleFunc("/events/{tag}", eventHandler.GetEventsByTag).Methods("GET")

	handlerWithAuth := authHandler.AuthMiddleware(whitelist, authHandler, r)
    handlerWithCORS := auth.CORSMiddleware(handlerWithAuth)
    handler := auth.LoggingMiddleware(handlerWithCORS)

	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
