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

func TestEventHandler_GetEventsByCategory(t *testing.T) {
	t.Parallel()

	getEventsByCategory := &pb.GetEventsByCategoryRequest{
		CategoryID: 1,
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
				req := httptest.NewRequest(http.MethodGet, "/events/category", nil)
				req = mux.SetURLVars(req, map[string]string{"category": "1"})
				return req
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

				serviceMock.EXPECT().GetEventsByCategory(gomock.Any(), getEventsByCategory).Return(events, nil)

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
			name: "No id",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/events/category", nil)
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventServiceClient(ctrl)

				return &EventHandler{
					EventService: serviceMock,
					logger:       logger,
				}
			},
			wantCode: http.StatusBadRequest,
			wantBody: &GetEventsResponse{},
		},
		{
			name: "Internal error",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/events/category", nil)
				req = mux.SetURLVars(req, map[string]string{"category": "1"})
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventServiceClient(ctrl)
				serviceMock.EXPECT().GetEventsByCategory(gomock.Any(), getEventsByCategory).Return(nil, status.Error(codes.NotFound, grpc.ErrInternal))

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
			tt.setupFunc(ctrl).GetEventsByCategory(recorder, tt.req)

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
