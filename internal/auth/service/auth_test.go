package service

import (
	"context"
	"testing"

	"kudago/internal/auth/service/mocks"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_Register(t *testing.T) {
	t.Parallel()

	user := models.User{
		ID:       0,
		Username: "test",
		Password: "test",
		Email:    "test",
	}

	type expected struct {
		user models.User
		err  error
	}

	tests := []struct {
		name      string
		input     models.User
		setupFunc func(ctrl *gomock.Controller) *service
		expected  expected
	}{
		{
			name:  "success register",
			input: user,
			setupFunc: func(ctrl *gomock.Controller) *service {
				mockUserDB := mocks.NewMockUserDB(ctrl)

				mockUserDB.EXPECT().
					UserExists(context.Background(), user).
					Return(false, nil)

				mockUserDB.EXPECT().
					CreateUser(context.Background(), user).
					Return(user, nil)
				return NewService(mockUserDB)
			},
			expected: expected{
				user: user,
				err:  nil,
			},
		},
		{
			name:  "username is taken",
			input: user,
			setupFunc: func(ctrl *gomock.Controller) *service {
				mockUserDB := mocks.NewMockUserDB(ctrl)

				mockUserDB.EXPECT().
					UserExists(context.Background(), user).
					Return(true, nil)

				return NewService(mockUserDB)
			},
			expected: expected{
				user: models.User{},
				err:  models.ErrEmailIsUsed,
			},
		},
		{
			name:  "internal error",
			input: user,
			setupFunc: func(ctrl *gomock.Controller) *service {
				mockUserDB := mocks.NewMockUserDB(ctrl)

				mockUserDB.EXPECT().
					UserExists(context.Background(), user).
					Return(false, models.ErrInternal)

				return NewService(mockUserDB)
			},
			expected: expected{
				user: models.User{},
				err:  models.ErrInternal,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			actual, err := tt.setupFunc(ctrl).Register(context.Background(), tt.input)

			assert.Equal(t, tt.expected.user, actual)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}
