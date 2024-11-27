package main

import (
	"log"
	"net"
	"net/http"

	"kudago/cmd/auth/config"
	proto "kudago/internal/auth/api"
	grpcAuth "kudago/internal/auth/grpc"
	authService "kudago/internal/auth/service"
	"kudago/internal/interceptors"
	"kudago/internal/logger"
	"kudago/internal/metrics"
	"kudago/internal/repository/postgres"
	userRepository "kudago/internal/repository/postgres/users"
	sessionRepository "kudago/internal/repository/redis/session"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"google.golang.org/grpc"
)

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

	pool, err := postgres.InitPostgres(conf.PostgresConfig, appLogger)
	if err != nil {
		log.Fatalf("Failed to connect to the postgres database", err)
	}
	defer pool.Close()

	listener, err := net.Listen("tcp", conf.ServiceAddr)
	if err != nil {
		log.Fatalf("Не удалось запустить gRPC-сервер auth: %v", err)
	}

	userDB := userRepository.NewDB(pool)
	sessionDB := sessionRepository.NewDB(&conf.RedisConfig)

	authService := authService.NewService(userDB)
	authServer := grpcAuth.NewServerAPI(&authService, sessionDB, appLogger)
	metrics.InitMetrics()

	grpc_prometheus.EnableHandlingTimeHistogram()
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.MetricsUnaryInterceptor("auth_service"),
			interceptors.PanicRecoveryInterceptor,
		),
	)

	proto.RegisterAuthServiceServer(grpcServer, authServer)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		metricsAddr := ":9091"
		log.Printf("Метрики доступны на %s/metrics", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, nil); err != nil {
			log.Fatalf("Не удалось запустить HTTP-сервер для метрик: %v", err)
		}
	}()

	log.Printf("gRPC сервер запущен на %s", conf.ServiceAddr)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска gRPC-сервера auth: %v", err)
	}
}
