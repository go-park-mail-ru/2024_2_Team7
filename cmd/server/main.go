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

	authHandler := &auth.Handler{
		UserDB:    *userDB,
		SessionDb: *sessionDB,
	}

	eventHandler := &events.Handler{
		EventDB: *eventsDB,
	}

	whitelist := []string{
		"/login",
		"/register",
		"/events",
	}

	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	r.HandleFunc("/session", authHandler.CheckSession).Methods("GET")
	r.HandleFunc("/events", eventHandler.GetAllEvents).Methods("GET")
	r.HandleFunc("/events/{tag}", eventHandler.GetEventsByTag).Methods("GET")

	handlerWithCORS := auth.CORSMiddleware(r)

	handler := authHandler.AuthMiddleware(whitelist, *authHandler, handlerWithCORS)

	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
