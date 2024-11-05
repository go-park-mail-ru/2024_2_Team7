package events

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httpErrors "kudago/internal/http/errors"
	"kudago/internal/http/events/mocks"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEventHandler_GetCategories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupFunc func(ctrl *gomock.Controller) *EventHandler
		wantCode  int
		wantBody  interface{}
	}{
		{
			name: "Успешное получение категорий",
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventsGetter(ctrl)
				serviceMock.EXPECT().GetCategories(gomock.Any()).Return([]models.Category{
					{Name: "Music"}, {Name: "Theatre"},
				}, nil)

				return &EventHandler{
					getter: serviceMock,
				}
			},
			wantCode: http.StatusOK,
			wantBody: []models.Category{{Name: "Music"}, {Name: "Theatre"}},
		},
		{
			name: "Внутренняя ошибка сервера",
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventsGetter(ctrl)
				serviceMock.EXPECT().GetCategories(gomock.Any()).Return(nil, models.ErrInternal)

				return &EventHandler{
					getter: serviceMock,
				}
			},
			wantCode: http.StatusInternalServerError,
			wantBody: httpErrors.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			handler := tt.setupFunc(ctrl)
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/categories", nil)

			handler.GetCategories(recorder, req)

			assert.Equal(t, tt.wantCode, recorder.Code)
		})
	}
}
