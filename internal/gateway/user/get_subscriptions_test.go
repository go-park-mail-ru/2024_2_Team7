package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"kudago/internal/gateway/user/mocks"
	"kudago/internal/logger"
	pb "kudago/internal/user/api"
	"kudago/internal/user/grpc"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUserHandler_GetSubscriptions(t *testing.T) {
	t.Parallel()

	getSubscriptions := &pb.GetSubscriptionsRequest{
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
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *UserHandlers {
				serviceMock := mocks.NewMockUserServiceClient(ctrl)
				users := &pb.GetSubscriptionsResponse{
					Users: []*pb.User{
						{
							ID:       1,
							Username: "user1",
						},
					},
				}

				serviceMock.EXPECT().GetSubscriptions(gomock.Any(), getSubscriptions).Return(users, nil)

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
			name: "No id",
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
			wantCode: http.StatusBadRequest,
			wantBody: &GetUsersResponse{},
		},
		{
			name: "Not found",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/profile/subscribe", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "1"})

				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *UserHandlers {
				serviceMock := mocks.NewMockUserServiceClient(ctrl)
				serviceMock.EXPECT().GetSubscriptions(gomock.Any(), getSubscriptions).Return(nil, status.Error(codes.NotFound, grpc.ErrUserNotFound))

				return &UserHandlers{
					UserService: serviceMock,
					logger:      logger,
				}
			},
			wantCode: http.StatusConflict,
			wantBody: &GetUsersResponse{},
		},
		{
			name: "Internal error",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/users", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *UserHandlers {
				serviceMock := mocks.NewMockUserServiceClient(ctrl)
				serviceMock.EXPECT().GetSubscriptions(gomock.Any(), getSubscriptions).Return(nil, status.Error(codes.Internal, grpc.ErrInternal))

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
			tt.setupFunc(ctrl).GetSubscriptions(recorder, tt.req)

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
