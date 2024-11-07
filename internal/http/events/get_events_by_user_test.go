package events

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"kudago/internal/http/events/mocks"
	"kudago/internal/http/utils"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEventHandler_GetEventsByUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		req       *http.Request
		setupFunc func(ctrl *gomock.Controller) *EventHandler
		wantCode  int
	}{
		{
			name: "Успешное получение событий пользователя",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/events/my", nil)
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventsGetter(ctrl)
				serviceMock.EXPECT().GetEventsByUser(gomock.Any(), 1, gomock.Any()).Return([]models.Event{}, nil)

				return &EventHandler{
					getter: serviceMock,
				}
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Нет активной сессии",
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/events/my", nil)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				return &EventHandler{
					service: mocks.NewMockEventService(ctrl),
				}
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "Внутренняя ошибка сервера",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/events/my", nil)
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventsGetter(ctrl)
				serviceMock.EXPECT().GetEventsByUser(gomock.Any(), 1, gomock.Any()).Return(nil, models.ErrInternal)

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
			tt.setupFunc(ctrl).GetEventsByUser(recorder, tt.req)

			assert.Equal(t, tt.wantCode, recorder.Code)
		})
	}
}
