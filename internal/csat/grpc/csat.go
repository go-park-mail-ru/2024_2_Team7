package grpc

import (
	"context"

	"kudago/internal/logger"
	"kudago/internal/models"
)

const (
	errInternal     = "internal error"
	errTestNotFound = "test not found"
)

type ServerAPI struct {
	// pb.UnimplementedCSATServiceServer
	service CSATService
	logger  *logger.Logger
}

type CSATService interface {
	GetTest(ctx context.Context, query string) (models.Test, error)
}

func NewServerAPI(service CSATService, logger *logger.Logger) *ServerAPI {
	return &ServerAPI{
		service: service,
		logger:  logger,
	}
}

/*func (s *ServerAPI) AddAnswers(ctx context.Context, req *pb.AddAnswersRequest) (*pb.Empty, error) {
}
*/
