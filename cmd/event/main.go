package main

import (
	"log"
	"net"

	"kudago/cmd/event/config"
	proto "kudago/internal/event/api"
	grpcEvent "kudago/internal/event/grpc"
	eventService "kudago/internal/event/service"
	"kudago/internal/logger"
	"kudago/internal/repository/postgres"
	eventRepository "kudago/internal/repository/postgres/events"

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

	eventDB := eventRepository.NewDB(pool)

	grpcServer := grpc.NewServer()
	eventService := eventService.NewService(eventDB)

	eventServer := grpcEvent.NewServerAPI(&eventService, eventDB, appLogger)
	proto.RegisterEventServiceServer(grpcServer, eventServer)

	log.Printf("gRPC сервер запущен на %s", conf.ServiceAddr)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска gRPC-сервера: %v", err)
	}
}
