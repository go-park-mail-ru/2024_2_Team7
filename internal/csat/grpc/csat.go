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
	errInternal      = "internal error"
	errTestNotFound  = "test not found"
	errAlreadyExists = "already answered"
)

type ServerAPI struct {
	pb.UnimplementedCSATServiceServer
	service CSATService
	logger  *logger.Logger
}

type CSATService interface {
	GetTest(ctx context.Context, query string) (models.Test, error)
	AddAnswers(ctx context.Context, answers []models.Answer, userID int) error
	GetStatistics(ctx context.Context) ([]models.Stats, error)
}

func NewServerAPI(service CSATService, logger *logger.Logger) *ServerAPI {
	return &ServerAPI{
		service: service,
		logger:  logger,
	}
}

func (s *ServerAPI) AddAnswers(ctx context.Context, req *pb.AddAnswersRequest) (*pb.Empty, error) {
	answers := make([]models.Answer, 0, len(req.Answers))
	for _, answer := range req.Answers {
		temp := models.Answer{
			QuestionID: int(answer.QuestionID),
			Value:      int(answer.Value),
		}

		answers = append(answers, temp)
	}

	err := s.service.AddAnswers(ctx, answers, int(req.UserID))
	if err != nil {
		if errors.Is(err, models.ErrForeignKeyViolation) {
			return nil, status.Error(codes.AlreadyExists, errAlreadyExists)
		}
		s.logger.Error(ctx, "add answers", err)
		return nil, status.Error(codes.Internal, errInternal)
	}
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

func (s *ServerAPI) GetStatistics(ctx context.Context, in *pb.Empty) (*pb.GetStatisticsResponse, error) {
	statistics, err := s.service.GetStatistics(ctx)
	if err != nil {
		s.logger.Error(ctx, "get stats", err)
		return nil, status.Error(codes.Internal, errInternal)
	}
	resp := toStatisticsPB(statistics)

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

func toStatisticsPB(statistics []models.Stats) *pb.GetStatisticsResponse {
	stats := make([]*pb.Stats, 0, len(statistics))

	for _, stat := range statistics {
		temp := &pb.Stats{
			ID:       int32(stat.ID),
			Question: stat.Question,
			Value:    int32(stat.Value),
		}
		stats = append(stats, temp)
	}

	return &pb.GetStatisticsResponse{
		Statistics: stats,
	}
}

func toQuestionPB(question models.Question) *pb.Question {
	return &pb.Question{
		Id:   int32(question.ID),
		Text: question.Text,
	}
}
