package main

import (
	"fmt"
	"log"
	"net/http"

	"kudago/cmd/server/config"
	_ "kudago/docs"
	authHandlers "kudago/internal/gateway/auth"
	csatHandlers "kudago/internal/gateway/csat"
	eventHandlers "kudago/internal/gateway/event"
	notificationHandlers "kudago/internal/gateway/notification"
	userHandlers "kudago/internal/gateway/user"

	"kudago/internal/logger"
	"kudago/internal/metrics"
	"kudago/internal/middleware"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

// swag init

// @title           Swagger Vihodnoy API
// @version         1.0
// @description     This is a description of the Vihodnoy server.
// @termsOfService  http://swagger.io/terms/

func main() {
	fmt.Println("hui")

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}

	appLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("Server failed to start logger: %v", err)
	}
	defer appLogger.Logger.Sync()

	authHandler, err := authHandlers.NewHandlers(conf.AuthServiceAddr, conf.ImageServiceAddr, appLogger)
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}

	userHandler, err := userHandlers.NewHandlers(conf.UserServiceAddr, appLogger)
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}

	eventHandler, err := eventHandlers.NewHandlers(conf.EventServiceAddr, conf.ImageServiceAddr, appLogger)
	if err != nil {
		log.Fatalf("Failed to connect to event service: %v", err)
	}

	csatHandler, err := csatHandlers.NewHandlers(conf.CSATServiceAddr, appLogger)
	if err != nil {
		log.Fatalf("Failed to connect to csat service: %v", err)
	}

	notificationHandler, err := notificationHandlers.NewHandlers(conf.NotificationServiceAddr, appLogger)
	if err != nil {
		log.Fatalf("Failed to connect to notification service: %v", err)
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
	r.HandleFunc("/profile", userHandler.UpdateUser).Methods(http.MethodPut)

	r.HandleFunc("/profile/subscribe/{id:[0-9]+}", userHandler.Subscribe).Methods(http.MethodPost)
	r.HandleFunc("/profile/subscribe/{id:[0-9]+}", userHandler.GetSubscriptions).Methods(http.MethodGet)
	r.HandleFunc("/profile/subscribe/{id:[0-9]+}", userHandler.Unsubscribe).Methods(http.MethodDelete)

	r.HandleFunc("/events/{id:[0-9]+}", eventHandler.GetEventByID).Methods(http.MethodGet)
	r.HandleFunc("/events/categories/{category:[0-9]+}", eventHandler.GetEventsByCategory).Methods(http.MethodGet)
	r.HandleFunc("/events", eventHandler.GetUpcomingEvents).Methods(http.MethodGet)
	r.HandleFunc("/events/past", eventHandler.GetPastEvents).Methods(http.MethodGet)
	r.HandleFunc("/events/subscription", eventHandler.GetSubscriptionEvents).Methods(http.MethodGet)

	r.HandleFunc("/test", csatHandler.GetTest).Methods(http.MethodGet)
	r.HandleFunc("/test", csatHandler.AddAnswers).Methods(http.MethodPost)
	r.HandleFunc("/stats", csatHandler.GetStatistics).Methods(http.MethodGet)

	r.HandleFunc("/categories", eventHandler.GetCategories).Methods(http.MethodGet)
	r.HandleFunc("/events/user/{id:[0-9]+}", eventHandler.GetEventsByUser).Methods(http.MethodGet)
	r.HandleFunc("/events/{id:[0-9]+}", eventHandler.UpdateEvent).Methods(http.MethodPut)
	r.HandleFunc("/events/{id:[0-9]+}", eventHandler.DeleteEvent).Methods(http.MethodDelete)
	r.HandleFunc("/events", eventHandler.AddEvent).Methods(http.MethodPost)
	r.HandleFunc("/events/search", eventHandler.SearchEvents).Methods(http.MethodGet)
	r.HandleFunc("/events/favorites", eventHandler.GetFavorites).Methods(http.MethodGet)
	r.HandleFunc("/events/favorites/{id:[0-9]+}", eventHandler.AddEventToFavorites).Methods(http.MethodPost)
	r.HandleFunc("/events/favorites/{id:[0-9]+}", eventHandler.DeleteEventFromFavorites).Methods(http.MethodDelete)

	r.HandleFunc("/notification", notificationHandler.GetNotifications).Methods(http.MethodGet)
	r.HandleFunc("/notification", notificationHandler.CreateNotification).Methods(http.MethodPost)
	fmt.Println("hui")
	handlerWithAuth := middleware.AuthMiddleware(authHandler.AuthService, r)
	handlerWithCORS := middleware.CORSMiddleware(handlerWithAuth)
	handlerWithLogging := middleware.LoggingMiddleware(handlerWithCORS, appLogger.Logger)
	handlerWithMetrics := middleware.MetricsMiddleware(handlerWithLogging, "server")
	handler := middleware.PanicMiddleware(handlerWithMetrics)

	r.Handle("/metrics", promhttp.Handler())
	metrics.InitMetrics()

	err = http.ListenAndServe(":"+conf.Port, handler)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
