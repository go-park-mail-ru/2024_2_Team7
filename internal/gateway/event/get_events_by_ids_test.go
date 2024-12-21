package events

import (
	"net/http"
	"net/http/httptest"
	"testing"

	pb "kudago/internal/event/api"
	"kudago/internal/gateway/event/mocks"
	"kudago/internal/logger"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestEventHandler_GetEventsByIDs(t *testing.T) {
	t.Parallel()

	getEventByIDRequest := &pb.GetEventsByIDsRequest{
		IDs: []int32{1},
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
			name: "Успешное получение событий",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/events/1", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventServiceClient(ctrl)
				event := &pb.Events{
					Events: []*pb.Event{
						{
							ID:          1,
							Title:       "user1",
							Description: "user1@mail.ru",
						},
					},
				}

				serviceMock.EXPECT().GetEventsByIDs(gomock.Any(), getEventByIDRequest).Return(event, nil)

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			tt.setupFunc(ctrl).getEventsByIDs(tt.req.Context(), []int{1})

			assert.Equal(t, tt.wantCode, recorder.Code)
		})
	}
}
