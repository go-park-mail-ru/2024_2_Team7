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
	sessionRepository "kudago/internal/repository/redis/session"

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

	authHandler := auth.NewAuthHandler(&authService, appLogger)
	eventHandler := events.NewEventHandler(&eventService, eventDB, appLogger)

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/register", authHandler.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/logout", authHandler.Logout).Methods(http.MethodPost)
	r.HandleFunc("/session", authHandler.CheckSession).Methods(http.MethodGet)

	r.HandleFunc("/profile", authHandler.Profile).Methods(http.MethodGet)
	r.HandleFunc("/profile", authHandler.UpdateUser).Methods(http.MethodPut)

	r.HandleFunc("/users/subscribe/{id:[0-9]+}", authHandler.Subscribe).Methods(http.MethodPost)
	r.HandleFunc("/users/subscribe/{id:[0-9]+}", authHandler.Unsubscribe).Methods(http.MethodDelete)

	r.HandleFunc("/events/{id:[0-9]+}", eventHandler.GetEventByID).Methods(http.MethodGet)
	r.HandleFunc("/events/categories/{category}", eventHandler.GetEventsByCategory).Methods(http.MethodGet)
	r.HandleFunc("/events", eventHandler.GetUpcomingEvents).Methods(http.MethodGet)
	r.HandleFunc("/events/past", eventHandler.GetPastEvents).Methods(http.MethodGet)
	r.HandleFunc("/events/subscription", eventHandler.GetSubscriptionEvents).Methods(http.MethodGet)

	r.HandleFunc("/categories", eventHandler.GetCategories).Methods(http.MethodGet)
	r.HandleFunc("/events/my", eventHandler.GetEventsByUser).Methods(http.MethodGet)
	r.HandleFunc("/events/{id:[0-9]+}", eventHandler.UpdateEvent).Methods(http.MethodPut)
	r.HandleFunc("/events/{id:[0-9]+}", eventHandler.DeleteEvent).Methods(http.MethodDelete)
	r.HandleFunc("/events", eventHandler.AddEvent).Methods(http.MethodPost)
	r.HandleFunc("/events/search", eventHandler.SearchEvents).Methods(http.MethodGet)
	r.HandleFunc("/events/favorites", eventHandler.GetFavorites).Methods(http.MethodGet)
	r.HandleFunc("/events/favorites/{id:[0-9]+}", eventHandler.AddEventToFavorites).Methods(http.MethodPost)
	r.HandleFunc("/events/favorites/{id:[0-9]+}", eventHandler.DeleteEventFromFavorites).Methods(http.MethodDelete)

	whitelist := []string{
		"/login",
		"/register",
		"/events",
		"/static",
		"/session",
		"/logout",
		"/docs",
		"/categories",
		"/swagger",
	}

	handlerWithAuth := middleware.AuthMiddleware(whitelist, sessionDB, r)
	handlerWithCORS := middleware.CORSMiddleware(handlerWithAuth)
	handlerWithLogging := middleware.LoggingMiddleware(handlerWithCORS, appLogger.Logger)
	handler := middleware.PanicMiddleware(handlerWithLogging)

	err = http.ListenAndServe(":"+conf.Port, handler)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
