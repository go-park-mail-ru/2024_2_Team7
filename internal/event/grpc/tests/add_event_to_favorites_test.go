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

func TestEventGRPC_AddEventToFavorites(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		req         *pb.FavoriteEvent
		setupFunc   func(ctrl *gomock.Controller) *event.ServerAPI
		expectedErr error
	}{
		{
			name: "success add event to favorites",
			req: &pb.FavoriteEvent{
				EventID: 1,
				UserID:  1,
			},
			setupFunc: func(ctrl *gomock.Controller) *event.ServerAPI {
				mockEventService := mocks.NewMockEventService(ctrl)
				mockEventGetter := mocks.NewMockEventsGetter(ctrl)

				logger, _ := logger.NewLogger()

				mockEventService.EXPECT().
					AddEventToFavorites(context.Background(), models.FavoriteEvent{EventID: 1, UserID: 1}).
					Return(nil)
				return event.NewServerAPI(mockEventService, mockEventGetter, logger)
			},
			expectedErr: nil,
		},
		{
			name: "event not found",
			req: &pb.FavoriteEvent{
				EventID: 1,
				UserID:  1,
			},
			setupFunc: func(ctrl *gomock.Controller) *event.ServerAPI {
				mockEventService := mocks.NewMockEventService(ctrl)
				mockEventGetter := mocks.NewMockEventsGetter(ctrl)

				logger, _ := logger.NewLogger()

				mockEventService.EXPECT().
					AddEventToFavorites(context.Background(), models.FavoriteEvent{EventID: 1, UserID: 1}).
					Return(models.ErrForeignKeyViolation)
				return event.NewServerAPI(mockEventService, mockEventGetter, logger)
			},
			expectedErr: status.Error(codes.NotFound, event.ErrEventNotFound),
		},
		{
			name: "already added",
			req: &pb.FavoriteEvent{
				EventID: 1,
				UserID:  1,
			},
			setupFunc: func(ctrl *gomock.Controller) *event.ServerAPI {
				mockEventService := mocks.NewMockEventService(ctrl)
				mockEventGetter := mocks.NewMockEventsGetter(ctrl)

				logger, _ := logger.NewLogger()

				mockEventService.EXPECT().
					AddEventToFavorites(context.Background(), models.FavoriteEvent{EventID: 1, UserID: 1}).
					Return(models.ErrNothingToInsert)
				return event.NewServerAPI(mockEventService, mockEventGetter, logger)
			},
			expectedErr: status.Error(codes.AlreadyExists, event.ErrAlreadyInFavorites),
		},
		{
			name: "internal error",
			req: &pb.FavoriteEvent{
				EventID: 1,
				UserID:  1,
			},
			setupFunc: func(ctrl *gomock.Controller) *event.ServerAPI {
				mockEventService := mocks.NewMockEventService(ctrl)
				mockEventGetter := mocks.NewMockEventsGetter(ctrl)

				logger, _ := logger.NewLogger()

				mockEventService.EXPECT().
					AddEventToFavorites(context.Background(), models.FavoriteEvent{EventID: 1, UserID: 1}).
					Return(models.ErrInternal)
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

			_, err := tt.setupFunc(ctrl).AddEventToFavorites(context.Background(), tt.req)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
