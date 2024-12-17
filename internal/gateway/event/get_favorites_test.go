package events

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	pb "kudago/internal/event/api"
	"kudago/internal/gateway/event/mocks"
	"kudago/internal/gateway/utils"
	"kudago/internal/logger"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthHandler_GetFavorites(t *testing.T) {
	t.Parallel()

	getFavoritesRequest := &pb.GetFavoritesRequest{
		UserID: int32(1),
		Params: &pb.PaginationParams{Limit: 1, Offset: 1},
	}

	logger, _ := logger.NewLogger()

	tests := []struct {
		name      string
		req       *http.Request
		setupFunc func(ctrl *gomock.Controller) *EventHandler
		wantCode  int
		wantBody  *GetEventsResponse
	}{
		{
			name: "Успешная проверка сессии",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/session", nil)
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventServiceClient(ctrl)
					event := &pb.Events{
						Events: &pb.Event{

							ID:          1,
							Title:       "user1",
							Description: "user1@mail.ru",
						}
				}

				serviceMock.EXPECT().GetFavorites(gomock.Any(), getFavoritesRequest).Return(user, nil)

				return &EventHandler{
					EventService: serviceMock,
					logger:       logger,
				}
			},
			wantCode: http.StatusOK,
			wantBody: &GetEventsResponse{
				Events: []EventResponse{
					ID:       1,
					Username: "user1",
					Email:    "user1@mail.ru",
				},
			},
		},
		{
			name: "Нет активной сессии",
			req:  httptest.NewRequest(http.MethodGet, "/session", nil),
			setupFunc: func(ctrl *gomock.Controller) *AuthHandlers {
				serviceMock := mocks.NewMockAuthServiceClient(ctrl)
				return &AuthHandlers{
					AuthService: serviceMock,
					logger:      logger,
				}
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Пользователь не найден",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/session", nil)
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *AuthHandlers {
				serviceMock := mocks.NewMockAuthServiceClient(ctrl)

				serviceMock.EXPECT().GetUser(gomock.Any(), getUserRequest).Return(nil, status.Error(codes.NotFound, event.ErrUserNotFound))

				return &AuthHandlers{
					AuthService: serviceMock,
					logger:      logger,
				}
			},
			wantCode: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			tt.setupFunc(ctrl).CheckSession(recorder, tt.req)

			assert.Equal(t, tt.wantCode, recorder.Code)

			if tt.wantBody != nil {
				var resp AuthResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, &resp)
			}
		})
	}
}
