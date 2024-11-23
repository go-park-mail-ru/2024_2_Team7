package grpc

import (
	"context"
	"errors"

	pb "kudago/internal/csat/api"
	"kudago/internal/logger"
	"kudago/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	errInternal     = "internal error"
	errTestNotFound = "test not found"
)

type ServerAPI struct {
	pb.UnimplementedCSATServiceServer
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

func (s *ServerAPI) AddAnswers(ctx context.Context, req *pb.AddAnswersRequest) (*pb.Empty, error) {
	return nil, nil
}

func (s *ServerAPI) GetTest(ctx context.Context, in *pb.GetTestRequest) (*pb.GetTestResponse, error) {
	test, err := s.service.GetTest(ctx, in.Query)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, errTestNotFound)
		}
		s.logger.Error(ctx, "get test", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	resp := toTestPB(test)

	return resp, nil
}

func toTestPB(test models.Test) *pb.GetTestResponse {
	questions := make([]*pb.Question, 0, len(test.Questions))

	for _, question := range test.Questions {
		q := toQuestionPB(question)
		questions = append(questions, q)
	}

	return &pb.GetTestResponse{
		Id:        int32(test.ID),
		Title:     test.Title,
		Questions: questions,
	}
}

func toQuestionPB(question models.Question) *pb.Question {
	return &pb.Question{
		Id:   int32(question.ID),
		Text: question.Text,
	}
}
