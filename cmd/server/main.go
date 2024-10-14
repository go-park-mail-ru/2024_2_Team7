package main

import (
	"log"
	"net/http"

	"kudago/config"
	handler "kudago/internal/http"
	"kudago/internal/middleware"
	"kudago/internal/repository"
	"kudago/internal/service"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	port := config.LoadConfig()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Server failed to start logger: %v", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	userDB := repository.NewUserDB()
	sessionDB := repository.NewSessionDB()
	eventDB := repository.NewEventDB()

	authService := service.NewAuthService(userDB, sessionDB)
	eventService := service.NewEventService(eventDB)

	authHandler := handler.NewAuthHandler(&authService)
	eventHandler := handler.NewEventHandler(&eventService)

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	r.HandleFunc("/session", authHandler.CheckSession).Methods("GET")
	r.HandleFunc("/events", eventHandler.GetAllEvents).Methods("GET")
	r.HandleFunc("/events/{tag}", eventHandler.GetEventsByTag).Methods("GET")

	whitelist := []string{
		"/login",
		"/register",
		"/events",
		"/static",
		"/session",
		"/logout",
	}

	handlerWithAuth := middleware.AuthMiddleware(whitelist, authHandler, r)
	handlerWithCORS := middleware.CORSMiddleware(handlerWithAuth)
	handlerWithLogging := middleware.LoggingMiddleware(handlerWithCORS, sugar)
	handler := middleware.PanicMiddleware(handlerWithLogging)

	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
