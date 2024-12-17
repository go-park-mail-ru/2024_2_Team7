package events

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	pbEvent "kudago/internal/event/api"
	"kudago/internal/gateway/event/mocks"
	"kudago/internal/gateway/utils"
	"kudago/internal/logger"
	"kudago/internal/models"
	pb "kudago/internal/notification/api"
	"kudago/internal/notification/grpc"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestEventHandler_GetNotifications(t *testing.T) {
	t.Parallel()

	getNotificationsRequest := &pb.GetNotificationsRequest{
		UserID: int32(1),
	}

	logger, _ := logger.NewLogger()

	tests := []struct {
		name      string
		req       *http.Request
		setupFunc func(ctrl *gomock.Controller) *EventHandler
		wantCode  int
		wantBody  *GetNotificationsResponse
	}{
		{
			name: "Успешное получение уведомлений",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/notifications", nil)
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceNotificationMock := mocks.NewMockNotificationServiceClient(ctrl)
				serviceEventMock := mocks.NewMockEventServiceClient(ctrl)

				notifications := &pb.GetNotificationsResponse{
					Notifications: []*pb.Notification{
						{
							Id:      1,
							EventID: 1,
						},
					},
				}

				events := &pbEvent.Events{
					Events: []*pbEvent.Event{
						{
							ID:          1,
							Title:       "user1",
							Description: "user1@mail.ru",
						},
					},
				}

				serviceNotificationMock.EXPECT().GetNotifications(gomock.Any(), getNotificationsRequest).Return(notifications, nil)
				serviceEventMock.EXPECT().GetEventsByIDs(gomock.Any(), gomock.Any()).Return(events, nil)

				return &EventHandler{
					NotificationService: serviceNotificationMock,
					EventService:        serviceEventMock,
					logger:              logger,
				}
			},
			wantCode: http.StatusOK,
			wantBody: &GetNotificationsResponse{
				Notifications: []NotificationWithEvent{
					{
						Notification: models.Notification{
							ID:      1,
							EventID: 1,
						},
						Event: models.Event{
							ID:          1,
							Title:       "user1",
							Description: "user1@mail.ru",
						},
					},
				},
			},
		},
		{
			name: "Internal error",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/notifications", nil)
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockNotificationServiceClient(ctrl)
				serviceMock.EXPECT().GetNotifications(gomock.Any(), getNotificationsRequest).Return(nil, status.Error(codes.Internal, grpc.ErrInternal))

				return &EventHandler{
					NotificationService: serviceMock,
					logger:              logger,
				}
			},
			wantCode: http.StatusInternalServerError,
			wantBody: &GetNotificationsResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			tt.setupFunc(ctrl).GetNotifications(recorder, tt.req)

			assert.Equal(t, tt.wantCode, recorder.Code)

			if tt.wantBody != nil {
				var resp GetNotificationsResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, &resp)
			}
		})
	}
}
