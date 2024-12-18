package main

import (
	"log"
	"net"
	"net/http"

	"kudago/cmd/user/config"
	"kudago/internal/interceptors"
	"kudago/internal/logger"
	"kudago/internal/metrics"
	"kudago/internal/repository/postgres"
	proto "kudago/internal/user/api"
	grpcEvent "kudago/internal/user/grpc"
	userRepository "kudago/internal/user/repository"

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
		log.Fatalf("Failed to connect to the postgres database: %v", err)
	}
	defer pool.Close()

	listener, err := net.Listen("tcp", conf.ServiceAddr)
	if err != nil {
		log.Fatalf("Не удалось запустить gRPC-сервер user: %v", err)
	}

	userDB := userRepository.NewDB(pool)
	metrics.InitMetrics()

	grpc_prometheus.EnableHandlingTimeHistogram()
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.MetricsUnaryInterceptor("user_service"),
			interceptors.PanicRecoveryInterceptor,
		),
	)

	userServer := grpcEvent.NewServerAPI(userDB, appLogger)
	proto.RegisterUserServiceServer(grpcServer, userServer)

	grpc_prometheus.Register(grpcServer)

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
