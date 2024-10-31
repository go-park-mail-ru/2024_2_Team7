package main

import (
	"log"
	"net/http"

	"kudago/config"
	_ "kudago/docs"
	"kudago/internal/http/auth"
	"kudago/internal/http/events"
	"kudago/internal/logger"
	"kudago/internal/middleware"
	imageRepository "kudago/internal/repository/images"
	"kudago/internal/repository/postgres"
	eventRepository "kudago/internal/repository/postgres/events"
	userRepository "kudago/internal/repository/postgres/users"
	sessionRepository "kudago/internal/repository/session"

	authService "kudago/internal/service/auth"
	eventService "kudago/internal/service/events"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
)

// swag init

// @title           Swagger Vihodnoy API
// @version         1.0
// @description     This is a description of the Vihodnoy server.
// @termsOfService  http://swagger.io/terms/

func main() {
	conf, err := config.LoadConfig()

	appLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("Server failed to start logger: %v", err)
	}
	defer appLogger.Logger.Sync()

	pool, err := postgres.InitPostgres(conf.PostgresConfig, appLogger)
	if err != nil {
		log.Fatalf("Failed to connect to the postgres database")
	}
	defer pool.Close()

	userDB := userRepository.NewDB(pool)
	sessionDB := sessionRepository.NewDB(&conf.RedisConfig)
	eventDB := eventRepository.NewDB(pool)
	imageDB := imageRepository.NewDB(conf.ImageConfig)
	authService := authService.NewService(userDB, sessionDB, imageDB)
	eventService := eventService.NewService(eventDB, imageDB)

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
	r.HandleFunc("/profile", authHandler.UpdateUser).Methods("PUT")

	r.HandleFunc("/events/{id:[0-9]+}", eventHandler.GetEventByID).Methods("GET")
	r.HandleFunc("/events/tags", eventHandler.GetEventsByTags).Methods("GET")
	r.HandleFunc("/events/categories/{category}", eventHandler.GetEventsByCategory).Methods("GET")
	r.HandleFunc("/events", eventHandler.GetUpcomingEvents).Methods("GET")
	r.HandleFunc("/pastevents", eventHandler.GetPastEvents).Methods("GET")

	r.HandleFunc("/categories", eventHandler.GetCategories).Methods("GET")
	r.HandleFunc("/events/my", eventHandler.GetEventsByUser).Methods("GET")
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
		"/docs",
		"/categories",
	}

	handlerWithAuth := middleware.AuthMiddleware(whitelist, authHandler, r)
	handlerWithCORS := middleware.CORSMiddleware(handlerWithAuth)
	handlerWithLogging := middleware.LoggingMiddleware(handlerWithCORS, appLogger.Logger)
	handler := middleware.PanicMiddleware(handlerWithLogging)

	err = http.ListenAndServe(":"+conf.Port, handler)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
