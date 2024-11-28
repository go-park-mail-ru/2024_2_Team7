package grpc

import (
	"context"
	"testing"

	pb "kudago/internal/csat/api"
	"kudago/internal/csat/grpc/mocks"
	"kudago/internal/logger"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCSATGRPC_GetStatistics(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		req         *pb.Empty
		setupFunc   func(ctrl *gomock.Controller) *ServerAPI
		expectedErr error
	}{
		{
			name: "success get statistics",
			req:  &pb.Empty{},
			setupFunc: func(ctrl *gomock.Controller) *ServerAPI {
				mockCSATService := mocks.NewMockCSATService(ctrl)
				logger, _ := logger.NewLogger()

				mockCSATService.EXPECT().
					GetStatistics(context.Background()).
					Return(nil, nil)
				return NewServerAPI(mockCSATService, logger)
			},
			expectedErr: nil,
		},
		{
			name: "internal error",
			req:  &pb.Empty{},
			setupFunc: func(ctrl *gomock.Controller) *ServerAPI {
				mockCSATService := mocks.NewMockCSATService(ctrl)
				logger, _ := logger.NewLogger()

				mockCSATService.EXPECT().
					GetStatistics(context.Background()).
					Return(nil, models.ErrInternal)
				return NewServerAPI(mockCSATService, logger)
			},
			expectedErr: status.Error(codes.Internal, errInternal),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			_, err := tt.setupFunc(ctrl).GetStatistics(context.Background(), tt.req)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestCSATGRPC_GetTest(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		req         *pb.GetTestRequest
		setupFunc   func(ctrl *gomock.Controller) *ServerAPI
		expectedRes *pb.GetTestResponse
		expectedErr error
	}{
		{
			name: "success get test",
			req:  &pb.GetTestRequest{Query: "test_query"},
			setupFunc: func(ctrl *gomock.Controller) *ServerAPI {
				mockCSATService := mocks.NewMockCSATService(ctrl)
				logger, _ := logger.NewLogger()

				mockCSATService.EXPECT().
					GetTest(context.Background(), "test_query").
					Return(models.Test{
						ID:    1,
						Title: "Test Title",
						Questions: []models.Question{
							{ID: 1, Text: "Question 1"},
							{ID: 2, Text: "Question 2"},
						},
					}, nil)

				return NewServerAPI(mockCSATService, logger)
			},
			expectedRes: &pb.GetTestResponse{
				Id:    1,
				Title: "Test Title",
				Questions: []*pb.Question{
					{Id: 1, Text: "Question 1"},
					{Id: 2, Text: "Question 2"},
				},
			},
			expectedErr: nil,
		},
		{
			name: "test not found",
			req:  &pb.GetTestRequest{Query: "invalid_query"},
			setupFunc: func(ctrl *gomock.Controller) *ServerAPI {
				mockCSATService := mocks.NewMockCSATService(ctrl)
				logger, _ := logger.NewLogger()

				mockCSATService.EXPECT().
					GetTest(context.Background(), "invalid_query").
					Return(models.Test{}, models.ErrNotFound)

				return NewServerAPI(mockCSATService, logger)
			},
			expectedRes: nil,
			expectedErr: status.Error(codes.NotFound, errTestNotFound),
		},
		{
			name: "internal error",
			req:  &pb.GetTestRequest{Query: "test_query"},
			setupFunc: func(ctrl *gomock.Controller) *ServerAPI {
				mockCSATService := mocks.NewMockCSATService(ctrl)
				logger, _ := logger.NewLogger()

				mockCSATService.EXPECT().
					GetTest(context.Background(), "test_query").
					Return(models.Test{}, models.ErrInternal)

				return NewServerAPI(mockCSATService, logger)
			},
			expectedRes: nil,
			expectedErr: status.Error(codes.Internal, errInternal),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := tt.setupFunc(ctrl)

			resp, err := server.GetTest(context.Background(), tt.req)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedRes, resp)
		})
	}
}

func TestCSATGRPC_AddAnswers(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		req         *pb.AddAnswersRequest
		setupFunc   func(ctrl *gomock.Controller) *ServerAPI
		expectedErr error
	}{
		{
			name: "success add answers",
			req: &pb.AddAnswersRequest{
				UserID: 1,
				Answers: []*pb.Answer{
					{QuestionID: 1, Value: 5},
					{QuestionID: 2, Value: 3},
				},
			},
			setupFunc: func(ctrl *gomock.Controller) *ServerAPI {
				mockCSATService := mocks.NewMockCSATService(ctrl)
				logger, _ := logger.NewLogger()

				mockCSATService.EXPECT().
					AddAnswers(context.Background(), []models.Answer{
						{QuestionID: 1, Value: 5},
						{QuestionID: 2, Value: 3},
					}, 1).
					Return(nil)

				return NewServerAPI(mockCSATService, logger)
			},
			expectedErr: nil,
		},
		{
			name: "already exists error",
			req: &pb.AddAnswersRequest{
				UserID: 1,
				Answers: []*pb.Answer{
					{QuestionID: 1, Value: 5},
				},
			},
			setupFunc: func(ctrl *gomock.Controller) *ServerAPI {
				mockCSATService := mocks.NewMockCSATService(ctrl)
				logger, _ := logger.NewLogger()

				mockCSATService.EXPECT().
					AddAnswers(context.Background(), []models.Answer{
						{QuestionID: 1, Value: 5},
					}, 1).
					Return(models.ErrForeignKeyViolation)

				return NewServerAPI(mockCSATService, logger)
			},
			expectedErr: status.Error(codes.AlreadyExists, errAlreadyExists),
		},
		{
			name: "internal error",
			req: &pb.AddAnswersRequest{
				UserID: 1,
				Answers: []*pb.Answer{
					{QuestionID: 1, Value: 5},
				},
			},
			setupFunc: func(ctrl *gomock.Controller) *ServerAPI {
				mockCSATService := mocks.NewMockCSATService(ctrl)
				logger, _ := logger.NewLogger()

				mockCSATService.EXPECT().
					AddAnswers(context.Background(), []models.Answer{
						{QuestionID: 1, Value: 5},
					}, 1).
					Return(models.ErrInternal)

				return NewServerAPI(mockCSATService, logger)
			},
			expectedErr: status.Error(codes.Internal, errInternal),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := tt.setupFunc(ctrl)

			_, err := server.AddAnswers(context.Background(), tt.req)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
