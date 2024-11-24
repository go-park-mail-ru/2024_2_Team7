package events

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"kudago/internal/http/events/mocks"
	"kudago/internal/logger"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEventHandler_GetPastEvents(t *testing.T) {
	t.Parallel()

	logger, _ := logger.NewLogger()

	tests := []struct {
		name      string
		req       *http.Request
		setupFunc func(ctrl *gomock.Controller) *EventHandler
		wantCode  int
	}{
		{
			name: "Успешное получение прошедших событий",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/events/past", nil)
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventsGetter(ctrl)
				serviceMock.EXPECT().GetPastEvents(gomock.Any(), gomock.Any()).Return([]models.Event{}, nil)

				return &EventHandler{
					getter: serviceMock,
					logger: logger,
				}
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Внутренняя ошибка сервера",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/events/past", nil)
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventsGetter(ctrl)
				serviceMock.EXPECT().GetPastEvents(gomock.Any(), gomock.Any()).Return(nil, models.ErrInternal)

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
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			tt.setupFunc(ctrl).GetPastEvents(recorder, tt.req)

			assert.Equal(t, tt.wantCode, recorder.Code)
		})
	}
}
