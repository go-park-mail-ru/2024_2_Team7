package events

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	pb "kudago/internal/event/api"
	"kudago/internal/event/grpc"
	"kudago/internal/gateway/event/mocks"
	"kudago/internal/gateway/utils"
	"kudago/internal/logger"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestEventHandler_GetSubscriptionEvents(t *testing.T) {
	t.Parallel()

	getSubscriptionEvents := &pb.GetSubscriptionsRequest{
		ID: 1,
		Params: &pb.PaginationParams{
			Limit:  30,
			Offset: 0,
		},
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
			name: "Успешное получение  событий",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/events", nil)
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventServiceClient(ctrl)
				events := &pb.Events{
					Events: []*pb.Event{
						{
							ID:          1,
							Title:       "user1",
							Description: "user1@mail.ru",
						},
					},
				}

				serviceMock.EXPECT().GetSubscriptionsEvents(gomock.Any(), getSubscriptionEvents).Return(events, nil)

				return &EventHandler{
					EventService: serviceMock,
					logger:       logger,
				}
			},
			wantCode: http.StatusOK,
			wantBody: &GetEventsResponse{
				Events: []EventResponse{
					{
						ID:          1,
						Title:       "user1",
						Description: "user1@mail.ru",
					},
				},
			},
		},
		{
			name: "No auth",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/events", nil)
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventServiceClient(ctrl)

				return &EventHandler{
					EventService: serviceMock,
					logger:       logger,
				}
			},
			wantCode: http.StatusForbidden,
			wantBody: &GetEventsResponse{},
		},
		{
			name: "Internal error",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/events", nil)
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventServiceClient(ctrl)
				serviceMock.EXPECT().GetSubscriptionsEvents(gomock.Any(), getSubscriptionEvents).Return(nil, status.Error(codes.NotFound, grpc.ErrInternal))

				return &EventHandler{
					EventService: serviceMock,
					logger:       logger,
				}
			},
			wantCode: http.StatusInternalServerError,
			wantBody: &GetEventsResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			tt.setupFunc(ctrl).GetSubscriptionEvents(recorder, tt.req)

			assert.Equal(t, tt.wantCode, recorder.Code)

			if tt.wantBody != nil {
				var resp GetEventsResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, &resp)
			}
		})
	}
}
