package main

import (
	"log"
	"net"

	"kudago/cmd/image/config"
	proto "kudago/internal/image/api"
	grpcImage "kudago/internal/image/grpc"
	"kudago/internal/logger"
	imageRepository "kudago/internal/repository/images"

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

	listener, err := net.Listen("tcp", conf.ServiceAddr)
	if err != nil {
		log.Fatalf("Не удалось запустить gRPC-сервер image: %v", err)
	}

	imageDB := imageRepository.NewDB(conf.ImageConfig)

	grpcServer := grpc.NewServer()
	imageServer := grpcImage.NewServerAPI(imageDB, appLogger)
	proto.RegisterImageServiceServer(grpcServer, imageServer)

	log.Printf("gRPC сервер запущен на %s", conf.ServiceAddr)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска gRPC-сервера image: %v", err)
	}
}
