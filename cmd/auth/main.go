package main

import (
	"log"
	"net"
	"os"

	"kudago/config"
	proto "kudago/internal/auth/api"
	grpcAuth "kudago/internal/auth/http"
	authService "kudago/internal/auth/service"
	"kudago/internal/logger"
	imageRepository "kudago/internal/repository/images"
	"kudago/internal/repository/postgres"
	userRepository "kudago/internal/repository/postgres/users"
	sessionRepository "kudago/internal/repository/redis/session"

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

	authServiceAddr := os.Getenv("AUTH_SERVICE_ADDR")
	if authServiceAddr == "" {
		log.Fatalf("AUTH_SERVICE_ADDR не задан")
	}

	listener, err := net.Listen("tcp", authServiceAddr)
	if err != nil {
		log.Fatalf("Не удалось запустить gRPC-сервер: %v", err)
	}

	userDB := userRepository.NewDB(pool)
	sessionDB := sessionRepository.NewDB(&conf.RedisConfig)
	imageDB := imageRepository.NewDB(conf.ImageConfig)

	authService := authService.NewService(userDB, imageDB)
	grpcServer := grpc.NewServer()
	authServer := grpcAuth.NewServerAPI(&authService, sessionDB, appLogger)
	proto.RegisterAuthServiceServer(grpcServer, authServer)

	log.Printf("gRPC сервер запущен на %s", authServiceAddr)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска gRPC-сервера: %v", err)
	}
}
