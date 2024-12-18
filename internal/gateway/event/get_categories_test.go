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
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestEventHandler_GetCategories(t *testing.T) {
	t.Parallel()

	logger, _ := logger.NewLogger()

	tests := []struct {
		name      string
		req       *http.Request
		setupFunc func(ctrl *gomock.Controller) *EventHandler
		wantCode  int
		wantBody  *GetCategoriesResponse
	}{
		{
			name: "Успешное получение категорий",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/categories", nil)
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventServiceClient(ctrl)
				categories := &pb.GetCategoriesResponse{
					Categories: []*pb.Category{
						{
							ID:   1,
							Name: "user1",
						},
					},
				}

				serviceMock.EXPECT().GetCategories(gomock.Any(), nil).Return(categories, nil)

				return &EventHandler{
					EventService: serviceMock,
					logger:       logger,
				}
			},
			wantCode: http.StatusOK,
			wantBody: &GetCategoriesResponse{
				Categories: []models.Category{
					{
						ID:   1,
						Name: "user1",
					},
				},
			},
		},
		{
			name: "Internal error",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/categories", nil)
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventServiceClient(ctrl)
				serviceMock.EXPECT().GetCategories(gomock.Any(), nil).Return(nil, status.Error(codes.NotFound, grpc.ErrInternal))

				return &EventHandler{
					EventService: serviceMock,
					logger:       logger,
				}
			},
			wantCode: http.StatusInternalServerError,
			wantBody: &GetCategoriesResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			tt.setupFunc(ctrl).GetCategories(recorder, tt.req)

			assert.Equal(t, tt.wantCode, recorder.Code)

			if tt.wantBody != nil {
				var resp GetCategoriesResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, &resp)
			}
		})
	}
}
