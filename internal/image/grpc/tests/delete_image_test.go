package grpc

import (
	"context"
	"testing"

	pb "kudago/internal/image/api"
	image "kudago/internal/image/grpc"
	"kudago/internal/image/grpc/tests/mocks"
	"kudago/internal/logger"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestImageGRPC_DeleteImage(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		req         *pb.DeleteRequest
		setupFunc   func(ctrl *gomock.Controller) *image.ServerAPI
		expectedErr error
	}{
		{
			name: "successful delete",
			req: &pb.DeleteRequest{
				FileUrl: "/test.png",
			},
			setupFunc: func(ctrl *gomock.Controller) *image.ServerAPI {
				mockService := mocks.NewMockImageService(ctrl)
				logger, _ := logger.NewLogger()

				mockService.EXPECT().
					DeleteImage(gomock.Any(), "/test.png").
					Return(nil)

				return image.NewServerAPI(mockService, logger)
			},
			expectedErr: nil,
		},
		{
			name: "internal error",
			req: &pb.DeleteRequest{
				FileUrl: "/test.png",
			},
			setupFunc: func(ctrl *gomock.Controller) *image.ServerAPI {
				mockService := mocks.NewMockImageService(ctrl)
				logger, _ := logger.NewLogger()

				mockService.EXPECT().
					DeleteImage(gomock.Any(), gomock.Any()).
					Return(models.ErrInternal)

				return image.NewServerAPI(mockService, logger)
			},
			expectedErr: status.Error(codes.Internal, image.ErrInternal),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := tt.setupFunc(ctrl)
			_, err := server.DeleteImage(context.Background(), tt.req)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
