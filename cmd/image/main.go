package main

import (
	"log"
	"net"
	"net/http"

	"kudago/cmd/image/config"
	proto "kudago/internal/image/api"
	grpcImage "kudago/internal/image/grpc"
	"kudago/internal/interceptors"
	"kudago/internal/logger"
	"kudago/internal/metrics"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	imageRepository "kudago/internal/repository/images"
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

	listener, err := net.Listen("tcp", conf.ServiceAddr)
	if err != nil {
		log.Fatalf("Не удалось запустить gRPC-сервер image: %v", err)
	}

	imageDB := imageRepository.NewDB(conf.ImageConfig)

	imageServer := grpcImage.NewServerAPI(imageDB, appLogger)
	metrics.InitMetrics()

	grpc_prometheus.EnableHandlingTimeHistogram()
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.MetricsUnaryInterceptor("auth_service"),
			interceptors.PanicRecoveryInterceptor,
		),
	)

	proto.RegisterImageServiceServer(grpcServer, imageServer)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		metricsAddr := ":9094"
		log.Printf("Метрики доступны на %s/metrics", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, nil); err != nil {
			log.Fatalf("Не удалось запустить HTTP-сервер для метрик: %v", err)
		}
	}()

	log.Printf("gRPC сервер запущен на %s", conf.ServiceAddr)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска gRPC-сервера image: %v", err)
	}
}
