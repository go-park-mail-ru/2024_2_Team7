package main

import (
	"log"
	"net"

	"kudago/cmd/auth/config"
	proto "kudago/internal/auth/api"
	grpcAuth "kudago/internal/auth/grpc"
	authService "kudago/internal/auth/service"
	"kudago/internal/logger"
	"kudago/internal/repository/postgres"
	userRepository "kudago/internal/repository/postgres/users"
	sessionRepository "kudago/internal/repository/redis/session"

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
	grpcServer := grpc.NewServer()
	authServer := grpcAuth.NewServerAPI(&authService, sessionDB, appLogger)
	proto.RegisterAuthServiceServer(grpcServer, authServer)

	log.Printf("gRPC сервер запущен на %s", conf.ServiceAddr)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска gRPC-сервера auth: %v", err)
	}
}
