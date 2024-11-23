package events

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"kudago/internal/http/events/mocks"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestEventHandler_GetEventByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		req       *http.Request
		setupFunc func(ctrl *gomock.Controller) *EventHandler
		wantCode  int
	}{
		{
			name: "Успешное получение события",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, "/events/1", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventsGetter(ctrl)
				serviceMock.EXPECT().GetEventByID(gomock.Any(), 1).Return(models.Event{}, nil)

				return &EventHandler{
					getter: serviceMock,
				}
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Некорректный ID события",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, "/events/abc", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "abc"})
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				return &EventHandler{
					service: mocks.NewMockEventService(ctrl),
				}
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "Событие не найдено",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, "/events/1", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventsGetter(ctrl)
				serviceMock.EXPECT().GetEventByID(gomock.Any(), 1).Return(models.Event{}, models.ErrEventNotFound)

				return &EventHandler{
					getter: serviceMock,
				}
			},
			wantCode: http.StatusNoContent,
		},
		{
			name: "Внутренняя ошибка сервера",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, "/events/1", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventsGetter(ctrl)
				serviceMock.EXPECT().GetEventByID(gomock.Any(), 1).Return(models.Event{}, models.ErrInternal)

				return &EventHandler{
					getter: serviceMock,
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
			tt.setupFunc(ctrl).GetEventByID(recorder, tt.req)

			assert.Equal(t, tt.wantCode, recorder.Code)
		})
	}
}
