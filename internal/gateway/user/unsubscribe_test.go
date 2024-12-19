package handlers

import (
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
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUserHandler_Unsubscribe(t *testing.T) {
	t.Parallel()

	subscriptionRequest := &pb.Subscription{
		SubscriberID: 1,
		FollowsID:    2,
	}

	logger, _ := logger.NewLogger()

	tests := []struct {
		name      string
		req       *http.Request
		setupFunc func(ctrl *gomock.Controller) *UserHandlers
		wantCode  int
	}{
		{
			name: "Успешное получение",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, "/subscribe", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *UserHandlers {
				serviceMock := mocks.NewMockUserServiceClient(ctrl)

				serviceMock.EXPECT().Unsubscribe(gomock.Any(), subscriptionRequest).Return(nil, nil)

				return &UserHandlers{
					UserService: serviceMock,
					logger:      logger,
				}
			},
			wantCode: http.StatusOK,
		},
		{
			name: "No id",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/profile/subscribe", nil)
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *UserHandlers {
				serviceMock := mocks.NewMockUserServiceClient(ctrl)

				return &UserHandlers{
					UserService: serviceMock,
					logger:      logger,
				}
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "Not found",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/profile", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *UserHandlers {
				serviceMock := mocks.NewMockUserServiceClient(ctrl)
				serviceMock.EXPECT().Unsubscribe(gomock.Any(), subscriptionRequest).Return(nil, status.Error(codes.NotFound, grpc.ErrUserNotFound))

				return &UserHandlers{
					UserService: serviceMock,
					logger:      logger,
				}
			},
			wantCode: http.StatusConflict,
		},
		{
			name: "Self subscription",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/profile", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "1"})

				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *UserHandlers {
				serviceMock := mocks.NewMockUserServiceClient(ctrl)

				return &UserHandlers{
					UserService: serviceMock,
					logger:      logger,
				}
			},
			wantCode: http.StatusConflict,
		},
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
				serviceMock := mocks.NewMockUserServiceClient(ctrl)
				serviceMock.EXPECT().Unsubscribe(gomock.Any(), subscriptionRequest).Return(nil, status.Error(codes.Internal, grpc.ErrInternal))

				return &UserHandlers{
					UserService: serviceMock,
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
			tt.setupFunc(ctrl).Unsubscribe(recorder, tt.req)

			assert.Equal(t, tt.wantCode, recorder.Code)
		})
	}
}
