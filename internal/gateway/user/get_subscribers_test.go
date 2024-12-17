package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"kudago/internal/gateway/user/mocks"
	"kudago/internal/gateway/utils"
	"kudago/internal/logger"
	"kudago/internal/models"
	pb "kudago/internal/user/api"
	"kudago/internal/user/grpc"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUserHandler_GetSubscribers(t *testing.T) {
	t.Parallel()

	getSubscribers := &pb.GetSubscribersRequest{
		ID: 1,
	}

	logger, _ := logger.NewLogger()

	tests := []struct {
		name      string
		req       *http.Request
		setupFunc func(ctrl *gomock.Controller) *UserHandlers
		wantCode  int
		wantBody  *GetUsersResponse
	}{
		{
			name: "Успешное получение  подписчиков",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/profile/subscribe", nil)
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *UserHandlers {
				serviceMock := mocks.NewMockUserServiceClient(ctrl)
				users := &pb.GetSubscribersResponse{
					Users: []*pb.User{
						{
							ID:       1,
							Username: "user1",
						},
					},
				}

				serviceMock.EXPECT().GetSubscribers(gomock.Any(), getSubscribers).Return(users, nil)

				return &UserHandlers{
					UserService: serviceMock,
					logger:      logger,
				}
			},
			wantCode: http.StatusOK,
			wantBody: &GetUsersResponse{
				Users: []UserResponse{
					{
						ID:       1,
						Username: "user1",
					},
				},
			},
		},
		{
			name: "No auth",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/profile/subscribe", nil)
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *UserHandlers {
				serviceMock := mocks.NewMockUserServiceClient(ctrl)

				return &UserHandlers{
					UserService: serviceMock,
					logger:      logger,
				}
			},
			wantCode: http.StatusForbidden,
			wantBody: &GetUsersResponse{},
		},
		{
			name: "Internal error",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/users", nil)
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *UserHandlers {
				serviceMock := mocks.NewMockUserServiceClient(ctrl)
				serviceMock.EXPECT().GetSubscribers(gomock.Any(), getSubscribers).Return(nil, status.Error(codes.NotFound, grpc.ErrInternal))

				return &UserHandlers{
					UserService: serviceMock,
					logger:      logger,
				}
			},
			wantCode: http.StatusInternalServerError,
			wantBody: &GetUsersResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			tt.setupFunc(ctrl).GetSubscribers(recorder, tt.req)

			assert.Equal(t, tt.wantCode, recorder.Code)

			if tt.wantBody != nil {
				var resp GetUsersResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, &resp)
			}
		})
	}
}
