package main

import (
	"log"
	"net"
	"os"

	"kudago/config"
	"kudago/internal/logger"
	"kudago/internal/repository/postgres"
	userRepository "kudago/internal/repository/postgres/users"
	proto "kudago/internal/user/api"
	grpcUser "kudago/internal/user/http"

	"google.golang.org/grpc"
)

func main() {
	conf, err := config.LoadConfig()

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

	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	if userServiceAddr == "" {
		log.Fatalf("USER_SERVICE_ADDR не задан")
	}

	listener, err := net.Listen("tcp", userServiceAddr)
	if err != nil {
		log.Fatalf("Не удалось запустить gRPC-сервер user: %v", err)
	}

	userDB := userRepository.NewDB(pool)

	grpcServer := grpc.NewServer()
	userServer := grpcUser.NewServerAPI(userDB, appLogger)
	proto.RegisterUserServiceServer(grpcServer, userServer)

	log.Printf("gRPC сервер запущен на %s", userServiceAddr)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска gRPC-сервера: %v", err)
	}
}
