package grpc

import (
	"context"
	"testing"

	pb "kudago/internal/event/api"
	"kudago/internal/event/grpc/tests/mocks"
	"kudago/internal/logger"
	"kudago/internal/models"

	event "kudago/internal/event/grpc"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestEventGRPC_GetCategories(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	categories := []models.Category{
		{
			ID:   1,
			Name: "test",
		},
	}

	tests := []struct {
		name         string
		req          *pb.Empty
		setupFunc    func(ctrl *gomock.Controller) *event.ServerAPI
		expectedResp *pb.GetCategoriesResponse
		expectedErr  error
	}{
		{
			name: "success get categories",
			req:  &pb.Empty{},
			setupFunc: func(ctrl *gomock.Controller) *event.ServerAPI {
				mockEventService := mocks.NewMockEventService(ctrl)
				mockEventGetter := mocks.NewMockEventsGetter(ctrl)

				logger, _ := logger.NewLogger()

				mockEventGetter.EXPECT().
					GetCategories(context.Background()).
					Return(categories, nil)
				return event.NewServerAPI(mockEventService, mockEventGetter, logger)
			},
			expectedResp: &pb.GetCategoriesResponse{
				Categories: []*pb.Category{
					{
						ID:   1,
						Name: "test",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "internal error",
			req:  &pb.Empty{},
			setupFunc: func(ctrl *gomock.Controller) *event.ServerAPI {
				mockEventService := mocks.NewMockEventService(ctrl)
				mockEventGetter := mocks.NewMockEventsGetter(ctrl)

				logger, _ := logger.NewLogger()

				mockEventGetter.EXPECT().
					GetCategories(context.Background()).
					Return(nil, models.ErrInternal)
				return event.NewServerAPI(mockEventService, mockEventGetter, logger)
			},
			expectedErr: status.Error(codes.Internal, event.ErrInternal),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			_, err := tt.setupFunc(ctrl).GetCategories(context.Background(), tt.req)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
