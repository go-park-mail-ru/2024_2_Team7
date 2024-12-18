package userRepository

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kudago/internal/models"
)

func TestUserRepository_Subscribe(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name          string
		susbscription models.Subscription
		mockSetup     func(m pgxmock.PgxConnIface)
		expectErr     error
	}{
		{
			name: "Успешное подписка",
			susbscription: models.Subscription{
				SubscriberID: 1,
				FollowsID:    1,
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`INSERT INTO SUBSCRIPTION (subscriber_id, follows_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`).
					WithArgs(1, 1).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			expectErr: nil,
		},
		{
			name: "Пользователь не найден",
			susbscription: models.Subscription{
				SubscriberID: 1,
				FollowsID:    2,
			}, mockSetup: func(m pgxmock.PgxConnIface) {
				var pgErr *pgconn.PgError
				m.ExpectExec(`INSERT INTO SUBSCRIPTION (subscriber_id, follows_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`).
					WithArgs(1, 2).
					WillReturnError(pgErr)
			},
			expectErr: models.ErrForeignKeyViolation,
		},
		{
			name: "Уже подписан",
			susbscription: models.Subscription{
				SubscriberID: 1,
				FollowsID:    2,
			}, mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`INSERT INTO SUBSCRIPTION (subscriber_id, follows_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`).
					WithArgs(1, 2).
					WillReturnResult(pgxmock.NewResult("INSERT", 0))
			},
			expectErr: models.ErrNothingToInsert,
		},
		{
			name: "ошибка",
			susbscription: models.Subscription{
				FollowsID:    3,
				SubscriberID: 1,
			}, mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`INSERT INTO SUBSCRIPTION (subscriber_id, follows_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`).
					WithArgs(1, 3).
					WillReturnError(fmt.Errorf("database error"))
			},
			expectErr: fmt.Errorf("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockConn, err := pgxmock.NewConn()
			require.NoError(t, err)
			defer mockConn.Close(ctx)

			tt.mockSetup(mockConn)

			db := NewDB(mockConn)

			err = db.Unsubscribe(context.Background(), tt.susbscription)
			if tt.expectErr != nil {
				assert.Error(t, tt.expectErr, err)
			}
		})
	}
}
