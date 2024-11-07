package events

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"kudago/internal/http/events/mocks"
	"kudago/internal/logger"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEventHandler_SearchEvents(t *testing.T) {
	t.Parallel()
	logger, _ := logger.NewLogger()

	tests := []struct {
		name      string
		req       *http.Request
		setupFunc func(ctrl *gomock.Controller) *EventHandler
		wantCode  int
	}{
		{
			name: "Успешный поиск событий",
			req: func() *http.Request {
				body := SearchRequest{
					Query:      "concert",
					EventStart: "2024-11-01",
					EventEnd:   "2024-12-01",
					Tags:       []string{"music", "live"},
					CategoryID: 1,
				}
				data, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodGet, "/events/search", bytes.NewBuffer(data))
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventService(ctrl)
				serviceMock.EXPECT().SearchEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Event{}, nil)

				return &EventHandler{
					service: serviceMock,
					logger:  logger,
				}
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Неверный формат запроса",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/events", bytes.NewBufferString("invalid data"))
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				return &EventHandler{
					service: mocks.NewMockEventService(ctrl),
					logger:  logger,
				}
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "Внутренняя ошибка сервера",
			req: func() *http.Request {
				body := SearchRequest{
					Query:      "concert",
					EventStart: "2024-11-01",
					EventEnd:   "2024-12-01",
					Tags:       []string{"music", "live"},
					CategoryID: 1,
				}
				data, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodGet, "/events/search", bytes.NewBuffer(data))
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventService(ctrl)
				serviceMock.EXPECT().SearchEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, models.ErrInternal)

				return &EventHandler{
					service: serviceMock,
					logger:  logger,
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
			tt.setupFunc(ctrl).SearchEvents(recorder, tt.req)

			assert.Equal(t, tt.wantCode, recorder.Code)
		})
	}
}
