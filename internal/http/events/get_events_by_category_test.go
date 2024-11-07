package events

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"kudago/internal/http/events/mocks"
	"kudago/internal/logger"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestEventHandler_GetEventsByCategory(t *testing.T) {
	t.Parallel()
	logger, _ := logger.NewLogger()

	tests := []struct {
		name      string
		req       *http.Request
		setupFunc func(ctrl *gomock.Controller) *EventHandler
		wantCode  int
	}{
		{
			name: "Успешное получение событий",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/events/categories/1", nil)
				req = mux.SetURLVars(req, map[string]string{"category": "1"})
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventsGetter(ctrl)
				serviceMock.EXPECT().GetEventsByCategory(gomock.Any(), 1, gomock.Any()).Return([]models.Event{}, nil)

				return &EventHandler{
					getter: serviceMock,
					logger: logger,
				}
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Некорректный ID категории",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, "/events/abc", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "abc"})
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
				req := httptest.NewRequest(http.MethodGet, "/events/categories/1", nil)
				req = mux.SetURLVars(req, map[string]string{"category": "1"})
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventsGetter(ctrl)
				serviceMock.EXPECT().GetEventsByCategory(gomock.Any(), 1, gomock.Any()).Return(nil, models.ErrInternal)

				return &EventHandler{
					getter: serviceMock,
					logger: logger,
				}
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			tt.setupFunc(ctrl).GetEventsByCategory(recorder, tt.req)

			assert.Equal(t, tt.wantCode, recorder.Code)
		})
	}
}
