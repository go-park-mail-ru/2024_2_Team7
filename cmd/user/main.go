package main

import (
	"log"
	"net"

	"kudago/cmd/user/config"
	"kudago/internal/logger"
	"kudago/internal/repository/postgres"
	userRepository "kudago/internal/repository/postgres/users"
	proto "kudago/internal/user/api"
	grpcUser "kudago/internal/user/grpc"

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

	userDB := userRepository.NewDB(pool)

	grpcServer := grpc.NewServer()
	userServer := grpcUser.NewServerAPI(userDB, appLogger)
	proto.RegisterUserServiceServer(grpcServer, userServer)

	log.Printf("gRPC сервер запущен на %s", conf.ServiceAddr)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска gRPC-сервера: %v", err)
	}
}
