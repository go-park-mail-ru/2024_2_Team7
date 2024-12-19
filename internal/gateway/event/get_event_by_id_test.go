package events

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	pb "kudago/internal/event/api"
	"kudago/internal/event/grpc"
	"kudago/internal/gateway/event/mocks"
	"kudago/internal/logger"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestEventHandler_GetEventsByID(t *testing.T) {
	t.Parallel()

	getEventByIDRequest := &pb.GetEventByIDRequest{
		ID: 1,
	}

	logger, _ := logger.NewLogger()

	tests := []struct {
		name      string
		req       *http.Request
		setupFunc func(ctrl *gomock.Controller) *EventHandler
		wantCode  int
		wantBody  *EventResponse
	}{
		{
			name: "Успешное получение событий",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/events/1", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventServiceClient(ctrl)
				event := &pb.Event{
					ID:          1,
					Title:       "user1",
					Description: "user1@mail.ru",
				}

				serviceMock.EXPECT().GetEventByID(gomock.Any(), getEventByIDRequest).Return(event, nil)

				return &EventHandler{
					EventService: serviceMock,
					logger:       logger,
				}
			},
			wantCode: http.StatusOK,
			wantBody: &EventResponse{
				ID:          1,
				Title:       "user1",
				Description: "user1@mail.ru",
			},
		},
		{
			name: "Not found",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/events/1", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventServiceClient(ctrl)
				serviceMock.EXPECT().GetEventByID(gomock.Any(), getEventByIDRequest).Return(nil, status.Error(codes.NotFound, grpc.ErrEventNotFound))

				return &EventHandler{
					EventService: serviceMock,
					logger:       logger,
				}
			},
			wantCode: http.StatusNotFound,
			wantBody: &EventResponse{},
		},
		{
			name: "Internal error",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/events/1", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventServiceClient(ctrl)
				serviceMock.EXPECT().GetEventByID(gomock.Any(), getEventByIDRequest).Return(nil, status.Error(codes.Internal, grpc.ErrInternal))

				return &EventHandler{
					EventService: serviceMock,
					logger:       logger,
				}
			},
			wantCode: http.StatusInternalServerError,
			wantBody: &EventResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			tt.setupFunc(ctrl).GetEventByID(recorder, tt.req)

			assert.Equal(t, tt.wantCode, recorder.Code)

			if tt.wantBody != nil {
				var resp EventResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, &resp)
			}
		})
	}
}
