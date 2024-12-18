package sessionRepository

import (
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"

	"kudago/internal/models"
)

func TestSessionDB_CheckSession(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name          string
		token         string
		mockSetup     func(mock redismock.ClientMock)
		expectedError error
	}{
		{
			name:  "Сессия найдена",
			token: "valid-token",
			mockSetup: func(mock redismock.ClientMock) {
				mock.ExpectGet("valid-token").SetVal("1")
			},
			expectedError: nil,
		},
		{
			name:  "Сессия не найдена",
			token: "missing-token",
			mockSetup: func(mock redismock.ClientMock) {
				mock.ExpectGet("missing-token").RedisNil()
			},
			expectedError: models.ErrUserNotFound,
		},
		{
			name:  "Ошибка Redis",
			token: "error-token",
			mockSetup: func(mock redismock.ClientMock) {
				mock.ExpectGet("error-token").SetErr(errors.New("redis error"))
			},
			expectedError: errors.New("redis error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRedis, mock := redismock.NewClientMock()
			tt.mockSetup(mock)

			db := &SessionDB{client: mockRedis}
			_, err := db.CheckSession(ctx, tt.token)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mock.ExpectationsWereMet()
		})
	}
}

func TestSessionDB_DeleteSession(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name          string
		token         string
		mockSetup     func(mock redismock.ClientMock)
		expectedError error
	}{
		{
			name:  "Успешное удаление сессии",
			token: "valid-token",
			mockSetup: func(mock redismock.ClientMock) {
				mock.ExpectDel("valid-token").SetVal(1)
			},
			expectedError: nil,
		},
		{
			name:  "Ошибка при удалении сессии",
			token: "error-token",
			mockSetup: func(mock redismock.ClientMock) {
				mock.ExpectDel("error-token").SetErr(errors.New("redis error"))
			},
			expectedError: errors.New("redis error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRedis, mock := redismock.NewClientMock()
			tt.mockSetup(mock)

			db := &SessionDB{client: mockRedis}
			err := db.DeleteSession(ctx, tt.token)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mock.ExpectationsWereMet()
		})
	}
}
