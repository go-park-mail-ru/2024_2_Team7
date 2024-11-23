package main

import (
	"log"
	"net/http"

	"kudago/cmd/server/config"
	_ "kudago/docs"
	authHandlers "kudago/internal/gateway/auth"
	csatHandlers "kudago/internal/gateway/csat"
	eventHandlers "kudago/internal/gateway/event"
	userHandlers "kudago/internal/gateway/user"

	"kudago/internal/logger"
	"kudago/internal/middleware"

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
	if err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}

	appLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("Server failed to start logger: %v", err)
	}
	defer appLogger.Logger.Sync()

	authHandler, err := authHandlers.NewAuthHandlers(conf.AuthServiceAddr, conf.ImageServiceAddr, appLogger)
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}

	userHandler, err := userHandlers.NewUserHandlers(conf.UserServiceAddr, appLogger)
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}

	eventHandler, err := eventHandlers.NewEventHandlers(conf.EventServiceAddr, conf.ImageServiceAddr, appLogger)
	if err != nil {
		log.Fatalf("Failed to connect to event service: %v", err)
	}

	csatHandler, err := csatHandlers.NewCSATHandlers(conf.CSATServiceAddr, appLogger)
	if err != nil {
		log.Fatalf("Failed to connect to csat service: %v", err)
	}

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/register", authHandler.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/logout", authHandler.Logout).Methods(http.MethodPost)
	r.HandleFunc("/session", authHandler.CheckSession).Methods(http.MethodGet)

	r.HandleFunc("/profile/{id:[0-9]+}", userHandler.Profile).Methods(http.MethodGet)
	// r.HandleFunc("/profile", userHandler.UpdateUser).Methods(http.MethodPut)

	r.HandleFunc("/profile/subscribe/{id:[0-9]+}", userHandler.Subscribe).Methods(http.MethodPost)
	r.HandleFunc("/profile/subscribe/{id:[0-9]+}", userHandler.GetSubscriptions).Methods(http.MethodGet)
	r.HandleFunc("/profile/subscribe/{id:[0-9]+}", userHandler.Unsubscribe).Methods(http.MethodDelete)

	r.HandleFunc("/events/{id:[0-9]+}", eventHandler.GetEventByID).Methods(http.MethodGet)
	r.HandleFunc("/events/categories/{category}", eventHandler.GetEventsByCategory).Methods(http.MethodGet)
	r.HandleFunc("/events", eventHandler.GetUpcomingEvents).Methods(http.MethodGet)
	r.HandleFunc("/events/past", eventHandler.GetPastEvents).Methods(http.MethodGet)
	r.HandleFunc("/events/subscription", eventHandler.GetSubscriptionEvents).Methods(http.MethodGet)

	r.HandleFunc("/test", csatHandler.GetTest).Methods(http.MethodGet)

	r.HandleFunc("/categories", eventHandler.GetCategories).Methods(http.MethodGet)
	r.HandleFunc("/events/user/{id:[0-9]+}", eventHandler.GetEventsByUser).Methods(http.MethodGet)
	r.HandleFunc("/events/{id:[0-9]+}", eventHandler.UpdateEvent).Methods(http.MethodPut)
	r.HandleFunc("/events/{id:[0-9]+}", eventHandler.DeleteEvent).Methods(http.MethodDelete)
	r.HandleFunc("/events", eventHandler.AddEvent).Methods(http.MethodPost)
	r.HandleFunc("/events/search", eventHandler.SearchEvents).Methods(http.MethodGet)
	r.HandleFunc("/events/favorites", eventHandler.GetFavorites).Methods(http.MethodGet)
	r.HandleFunc("/events/favorites/{id:[0-9]+}", eventHandler.AddEventToFavorites).Methods(http.MethodPost)
	r.HandleFunc("/events/favorites/{id:[0-9]+}", eventHandler.DeleteEventFromFavorites).Methods(http.MethodDelete)

	handlerWithAuth := middleware.AuthMiddleware(authHandler.AuthService, r)
	handlerWithCORS := middleware.CORSMiddleware(handlerWithAuth)
	handlerWithLogging := middleware.LoggingMiddleware(handlerWithCORS, appLogger.Logger)
	handler := middleware.PanicMiddleware(handlerWithLogging)

	err = http.ListenAndServe(":"+conf.Port, handler)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
