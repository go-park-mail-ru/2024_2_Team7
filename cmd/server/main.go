package main

import (
	"log"
	"net/http"

	"kudago/config"
	_ "kudago/docs"
	"kudago/internal/http/auth"
	"kudago/internal/http/events"
	"kudago/internal/middleware"
	eventRepository "kudago/internal/repository/events"
	sessionRepository "kudago/internal/repository/session"
	userRepository "kudago/internal/repository/users"
	authService "kudago/internal/service/auth"
	eventService "kudago/internal/service/events"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// swag init

// @title           Swagger Vihodnoy API
// @version         1.0
// @description     This is a description of the Vihodnoy server.
// @termsOfService  http://swagger.io/terms/

func main() {
	port := config.LoadConfig()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Server failed to start logger: %v", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	userDB := userRepository.NewDB()
	sessionDB := sessionRepository.NewDB()
	eventDB := eventRepository.NewDB()

	authService := authService.NewService(userDB, sessionDB)
	eventService := eventService.NewService(eventDB)

	authHandler := auth.NewAuthHandler(&authService)
	eventHandler := events.NewEventHandler(&eventService)

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	r.HandleFunc("/session", authHandler.CheckSession).Methods("GET")
	r.HandleFunc("/profile", authHandler.Profile).Methods("GET")

	r.HandleFunc("/events/{id:[0-9]+}", eventHandler.GetEventByID).Methods("GET")
	r.HandleFunc("/events/{tag}", eventHandler.GetEventsByTag).Methods("GET")
	r.HandleFunc("/events", eventHandler.GetAllEvents).Methods("GET")
	r.HandleFunc("/events/{id:[0-9]+}", eventHandler.UpdateEvent).Methods("PUT")
	r.HandleFunc("/events/{id:[0-9]+}", eventHandler.DeleteEvent).Methods("DELETE")
	r.HandleFunc("/events", eventHandler.AddEvent).Methods("POST")

	whitelist := []string{
		"/login",
		"/register",
		"/events",
		"/static",
		"/session",
		"/logout",
		"/swagger",
		"/docs",
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
