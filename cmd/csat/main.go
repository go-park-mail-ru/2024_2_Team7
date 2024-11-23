package main

import (
	"fmt"
	"log"
	"net"

	"kudago/cmd/csat/config"
	proto "kudago/internal/csat/api"
	grpcCSAT "kudago/internal/csat/grpc"
	csatRepository "kudago/internal/csat/repository"
	"kudago/internal/logger"
	"kudago/internal/repository/postgres"

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
		log.Fatalf("Не удалось запустить gRPC-сервер user: %v", err)
	}

	csatDB := csatRepository.NewDB(pool)

	grpcServer := grpc.NewServer()
	csatServer := grpcCSAT.NewServerAPI(csatDB, appLogger)
	proto.RegisterCSATServiceServer(grpcServer, csatServer)
	fmt.Println(conf.PostgresConfig)
	log.Printf("gRPC сервер запущен на %s", conf.ServiceAddr)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска gRPC-сервера: %v", err)
	}
}
