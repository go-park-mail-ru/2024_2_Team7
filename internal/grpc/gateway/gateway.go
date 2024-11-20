package gateway

import (
	"time"

	pb "kudago/internal/auth/api"

	"google.golang.org/grpc"
)

type Gateway struct {
	authClient pb.AuthServiceClient
}

func NewGateway(authServiceAddr string) (*Gateway, error) {
	conn, err := grpc.Dial(authServiceAddr, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(5*time.Second))
	if err != nil {
		return nil, err
	}

	return &Gateway{
		authClient: pb.NewAuthServiceClient(conn),
	}, nil
}
