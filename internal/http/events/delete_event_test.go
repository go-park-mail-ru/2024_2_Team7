package events

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"kudago/internal/http/events/mocks"
	"kudago/internal/http/utils"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestEventHandler_DeleteEvent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		req       *http.Request
		setupFunc func(ctrl *gomock.Controller) *EventHandler
		wantCode  int
	}{
		{
			name: "Успешное удаление события",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, "/events/1", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventService(ctrl)
				serviceMock.EXPECT().DeleteEvent(gomock.Any(), 1, 1).Return(nil)

				return &EventHandler{
					service: serviceMock,
				}
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Нет активной сессии",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, "/events/1", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				return &EventHandler{
					service: mocks.NewMockEventService(ctrl),
				}
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "Некорректный ID события",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, "/events/abc", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "abc"})

				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
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
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventService(ctrl)
				serviceMock.EXPECT().DeleteEvent(gomock.Any(), 1, 1).Return(models.ErrEventNotFound)

				return &EventHandler{
					service: serviceMock,
				}
			},
			wantCode: http.StatusNotFound,
		},
		{
			name: "Доступ запрещен",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, "/events/1", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "1"})

				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventService(ctrl)
				serviceMock.EXPECT().DeleteEvent(gomock.Any(), 1, 1).Return(models.ErrAccessDenied)

				return &EventHandler{
					service: serviceMock,
				}
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "Внутренняя ошибка сервера",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, "/events/1", nil)
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventService(ctrl)
				serviceMock.EXPECT().DeleteEvent(gomock.Any(), 1, 1).Return(models.ErrInternal)

				return &EventHandler{
					service: serviceMock,
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
			tt.setupFunc(ctrl).DeleteEvent(recorder, tt.req)

			assert.Equal(t, tt.wantCode, recorder.Code)
		})
	}
}
