package main

import (
	"log"
	"net"
	"net/http"

	"kudago/cmd/user/config"
	"kudago/internal/logger"
	"kudago/internal/metrics"
	"kudago/internal/repository/postgres"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"kudago/internal/interceptors"
	userRepository "kudago/internal/repository/postgres/users"
	proto "kudago/internal/user/api"
	grpcUser "kudago/internal/user/grpc"

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
		log.Fatalf("Failed to connect to the postgres database", err, conf.PostgresConfig)
	}
	defer pool.Close()

	listener, err := net.Listen("tcp", conf.ServiceAddr)
	if err != nil {
		log.Fatalf("Не удалось запустить gRPC-сервер user: %v", err)
	}

	userDB := userRepository.NewDB(pool)

	userServer := grpcUser.NewServerAPI(userDB, appLogger)

	metrics.InitMetrics()

	grpc_prometheus.EnableHandlingTimeHistogram()
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.MetricsUnaryInterceptor("user_service"),
			interceptors.PanicRecoveryInterceptor,
		),
	)

	proto.RegisterUserServiceServer(grpcServer, userServer)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		metricsAddr := ":9092"
		log.Printf("Метрики доступны на %s/metrics", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, nil); err != nil {
			log.Fatalf("Не удалось запустить HTTP-сервер для метрик: %v", err)
		}
	}()

	log.Printf("gRPC сервер запущен на %s", conf.ServiceAddr)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска gRPC-сервера: %v", err)
	}
}
