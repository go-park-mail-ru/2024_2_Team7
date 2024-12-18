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

func TestEventGRPC_GetPastEvents(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventsData := []models.Event{
		{
			ID:    1,
			Title: "test",
		},
	}

	tests := []struct {
		name         string
		req          *pb.PaginationParams
		setupFunc    func(ctrl *gomock.Controller) *event.ServerAPI
		expectedResp *pb.Events
		expectedErr  error
	}{
		{
			name: "success get past events",
			req: &pb.PaginationParams{
				Limit:  1,
				Offset: 0,
			},
			setupFunc: func(ctrl *gomock.Controller) *event.ServerAPI {
				mockEventService := mocks.NewMockEventService(ctrl)
				mockEventGetter := mocks.NewMockEventsGetter(ctrl)

				logger, _ := logger.NewLogger()

				mockEventGetter.EXPECT().
					GetPastEvents(context.Background(), models.PaginationParams{Limit: 1, Offset: 0}).
					Return(eventsData, nil)
				return event.NewServerAPI(mockEventService, mockEventGetter, logger)
			},
			expectedResp: &pb.Events{
				Events: []*pb.Event{
					{
						ID:    1,
						Title: "test",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "internal error",
			req: &pb.PaginationParams{
				Limit:  1,
				Offset: 0,
			},
			setupFunc: func(ctrl *gomock.Controller) *event.ServerAPI {
				mockEventService := mocks.NewMockEventService(ctrl)
				mockEventGetter := mocks.NewMockEventsGetter(ctrl)

				logger, _ := logger.NewLogger()

				mockEventGetter.EXPECT().
					GetPastEvents(context.Background(), models.PaginationParams{Limit: 1, Offset: 0}).
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

			_, err := tt.setupFunc(ctrl).GetPastEvents(context.Background(), tt.req)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
