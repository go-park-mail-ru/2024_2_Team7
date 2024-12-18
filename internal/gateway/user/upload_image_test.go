package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"kudago/internal/gateway/user/mocks"
	"kudago/internal/gateway/utils"
	pb "kudago/internal/image/api"
	"kudago/internal/logger"
	"kudago/internal/models"
	"kudago/internal/user/grpc"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUserHandler_UploadImage(t *testing.T) {
	t.Parallel()

	uploadRequest := &pb.UploadRequest{
		Filename:    "2",
	}

	logger, _ := logger.NewLogger()

	tests := []struct {
		name      string
		req       *http.Request
		setupFunc func(ctrl *gomock.Controller) *UserHandlers
		wantCode  int
	}{
		// {
		// 	name: "Успешное получение",
		// 	req: func() *http.Request {
		// 		req := httptest.NewRequest(http.MethodDelete, "/subscribe", nil)
		// 		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		// 		session := models.Session{UserID: 1, Token: "valid_token"}
		// 		ctx := utils.SetSessionInContext(req.Context(), session)
		// 		return req.WithContext(ctx)
		// 	}(),
		// 	setupFunc: func(ctrl *gomock.Controller) *UserHandlers {
		// 		serviceMock := mocks.NewMockImageServiceClient(ctrl)

		// 		serviceMock.EXPECT().UploadImage(gomock.Any(), uploadRequest).Return(&pb.UploadResponse{}, nil)

		// 		return &UserHandlers{
		// 			ImageService: serviceMock,
		// 			logger:      logger,
		// 		}
		// 	},
		// 	wantCode: http.StatusOK,
		// },
		{
			name: "Internal error",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/users", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *UserHandlers {
				serviceMock := mocks.NewMockImageServiceClient(ctrl)
				serviceMock.EXPECT().UploadImage(gomock.Any(), uploadRequest).Return(nil, status.Error(codes.Internal, grpc.ErrInternal))

				return &UserHandlers{
					ImageService: serviceMock,
					logger:      logger,
				}
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			tt.setupFunc(ctrl).uploadImage(context.Background(), uploadRequest, recorder)

			assert.Equal(t, tt.wantCode, recorder.Code)
		})
	}
}
