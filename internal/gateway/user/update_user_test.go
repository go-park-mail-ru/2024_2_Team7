package handlers

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"kudago/internal/gateway/user/mocks"
// 	"kudago/internal/gateway/utils"
// 	"kudago/internal/logger"
// 	"kudago/internal/models"
// 	pb "kudago/internal/user/api"
// 	"kudago/internal/user/grpc"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// func TestUserHandler_UpdateUser(t *testing.T) {
// 	t.Parallel()

// 	updateRequest := &pb.User{
// 		ID:       1,
// 		Username: "user1",
// 		Email:    "user@mail.com",
// 	}

// 	logger, _ := logger.NewLogger()

// 	tests := []struct {
// 		name      string
// 		req       *http.Request
// 		setupFunc func(ctrl *gomock.Controller) *UserHandlers
// 		wantCode  int
// 		wantBody  UserResponse
// 	}{
// 		{
// 			name: "Успешное получение",
// 			req: func() *http.Request {
// 				req := httptest.NewRequest(http.MethodPost, "/profile", nil)
// 				session := models.Session{UserID: 1, Token: "valid_token"}
// 				ctx := utils.SetSessionInContext(req.Context(), session)
// 				return req.WithContext(ctx)
// 			}(),
// 			setupFunc: func(ctrl *gomock.Controller) *UserHandlers {
// 				serviceMock := mocks.NewMockUserServiceClient(ctrl)
// 				user := &pb.User{
// 					ID:       1,
// 					Username: "user1",
// 				}
// 				serviceMock.EXPECT().UpdateUser(gomock.Any(), updateRequest).Return(user, nil)

// 				return &UserHandlers{
// 					UserService: serviceMock,
// 					logger:      logger,
// 				}
// 			},
// 			wantCode: http.StatusOK,
// 			wantBody: UserResponse{
// 				ID:       1,
// 				Username: "user1",
// 			},
// 		},

// 		{
// 			name: "Username is taken",
// 			req: func() *http.Request {
// 				req := httptest.NewRequest(http.MethodGet, "/profile", nil)
// 				session := models.Session{UserID: 1, Token: "valid_token"}
// 				ctx := utils.SetSessionInContext(req.Context(), session)
// 				return req.WithContext(ctx)
// 			}(),
// 			setupFunc: func(ctrl *gomock.Controller) *UserHandlers {
// 				serviceMock := mocks.NewMockUserServiceClient(ctrl)
// 				serviceMock.EXPECT().UpdateUser(gomock.Any(), updateRequest).Return(nil, status.Error(codes.AlreadyExists, grpc.ErrUsernameOrEmailIsTaken))

// 				return &UserHandlers{
// 					UserService: serviceMock,
// 					logger:      logger,
// 				}
// 			},
// 			wantCode: http.StatusConflict,
// 			wantBody: UserResponse{},
// 		},
// 		{
// 			name: "Internal error",
// 			req: func() *http.Request {
// 				req := httptest.NewRequest(http.MethodGet, "/users", nil)
// 				session := models.Session{UserID: 1, Token: "valid_token"}
// 				ctx := utils.SetSessionInContext(req.Context(), session)
// 				return req.WithContext(ctx)
// 			}(),
// 			setupFunc: func(ctrl *gomock.Controller) *UserHandlers {
// 				serviceMock := mocks.NewMockUserServiceClient(ctrl)
// 				serviceMock.EXPECT().UpdateUser(gomock.Any(), updateRequest).Return(nil, status.Error(codes.Internal, grpc.ErrInternal))

// 				return &UserHandlers{
// 					UserService: serviceMock,
// 					logger:      logger,
// 				}
// 			},
// 			wantCode: http.StatusInternalServerError,
// 			wantBody: UserResponse{},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			recorder := httptest.NewRecorder()
// 			tt.setupFunc(ctrl).UpdateUser(recorder, tt.req)

// 			assert.Equal(t, tt.wantCode, recorder.Code)
// 		})
// 	}
// }
