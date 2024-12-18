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

func TestImageGRPC_UploadImage(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		req         *pb.UploadRequest
		setupFunc   func(ctrl *gomock.Controller) *image.ServerAPI
		expectedRes *pb.UploadResponse
		expectedErr error
	}{
		{
			name: "successful upload",
			req: &pb.UploadRequest{
				Filename: "test.png",
				File:     []byte("image-data"),
			},
			setupFunc: func(ctrl *gomock.Controller) *image.ServerAPI {
				mockService := mocks.NewMockImageService(ctrl)
				logger, _ := logger.NewLogger()

				mockService.EXPECT().
					UploadImage(gomock.Any(), gomock.Any()).
					Return("/test.png", nil)

				return image.NewServerAPI(mockService, logger)
			},
			expectedRes: &pb.UploadResponse{
				FileUrl: "/test.png",
			},
			expectedErr: nil,
		},
		{
			name: "internal error",
			req: &pb.UploadRequest{
				Filename: "test.png",
				File:     []byte("image-data"),
			},
			setupFunc: func(ctrl *gomock.Controller) *image.ServerAPI {
				mockService := mocks.NewMockImageService(ctrl)
				logger, _ := logger.NewLogger()

				mockService.EXPECT().
					UploadImage(gomock.Any(), gomock.Any()).
					Return("", models.ErrInternal)

				return image.NewServerAPI(mockService, logger)
			},
			expectedRes: nil,
			expectedErr: status.Error(codes.Internal, image.ErrInternal),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := tt.setupFunc(ctrl)
			res, err := server.UploadImage(context.Background(), tt.req)

			assert.Equal(t, tt.expectedRes, res)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
