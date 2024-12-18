package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"kudago/internal/notification/repository"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kudago/internal/models"
)

func TestEventRepository_Notification(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		eventID   int
		mockSetup func(m pgxmock.PgxConnIface)
		expectErr bool
	}{
		{
			name:    "Успешное удаление",
			eventID: 1,
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`DELETE FROM NOTIFICATION WHERE id=\$1`).
					WithArgs(1).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			expectErr: false,
		},
		{
			name:    "Событие не найдено",
			eventID: 2,
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`DELETE FROM NOTIFICATION WHERE id=\$1`).
					WithArgs(2).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))
			},
			expectErr: false,
		},
		{
			name:    "ошибка",
			eventID: 3,
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`DELETE FROM NOTIFICATION WHERE id=\$1`).
					WithArgs(3).
					WillReturnError(fmt.Errorf("database error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockConn, err := pgxmock.NewConn()
			require.NoError(t, err)
			defer mockConn.Close(ctx)

			tt.mockSetup(mockConn)

			db := repository.NewDB(mockConn)

			err = db.DeleteNotification(context.Background(), tt.eventID)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), models.LevelDB)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNotificationRepository_CreateNotification(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Формат времени для парсинга строк в time.Time
	timeLayout := "2006-01-02 15:04:05"

	tests := []struct {
		name         string
		notification models.Notification
		mockSetup    func(m pgxmock.PgxConnIface)
		expectErr    bool
	}{
		{
			name: "Успешное создание уведомления",
			notification: models.Notification{
				UserID:   1,
				EventID:  1,
				Message:  "Event Reminder",
				NotifyAt: parseTime(t, "2024-12-18 10:00:00", timeLayout),
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`INSERT INTO NOTIFICATION`).
					WithArgs(1, 1, "Event Reminder", parseTime(t, "2024-12-18 10:00:00", timeLayout)).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			expectErr: false,
		},
		{
			name: "Ошибка при вставке уведомления",
			notification: models.Notification{
				UserID:   2,
				EventID:  2,
				Message:  "Another Event",
				NotifyAt: parseTime(t, "2024-12-19 12:00:00", timeLayout),
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`INSERT INTO NOTIFICATION`).
					WithArgs(2, 2, "Another Event", parseTime(t, "2024-12-19 12:00:00", timeLayout)).
					WillReturnError(fmt.Errorf("database error"))
			},
			expectErr: true,
		},
		{
			name: "Неверные данные",
			notification: models.Notification{
				UserID:   0,
				EventID:  3,
				Message:  "Invalid User",
				NotifyAt: parseTime(t, "2024-12-20 15:00:00", timeLayout),
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectExec(`INSERT INTO NOTIFICATION`).
					WithArgs(0, 3, "Invalid User", parseTime(t, "2024-12-20 15:00:00", timeLayout)).
					WillReturnError(fmt.Errorf("invalid user id"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockConn, err := pgxmock.NewConn()
			require.NoError(t, err)
			defer mockConn.Close(ctx)

			tt.mockSetup(mockConn)

			db := repository.NewDB(mockConn)

			err = db.CreateNotification(context.Background(), tt.notification)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), models.LevelDB)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNotificationRepository_CreateNotificationsByUserIDs(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Формат времени для парсинга строк в time.Time
	timeLayout := "2006-01-02 15:04:05"

	tests := []struct {
		name         string
		ids          []int
		notification models.Notification
		mockSetup    func(m pgxmock.PgxConnIface)
		expectErr    bool
	}{
		{
			name: "Успешное создание уведомлений для нескольких пользователей",
			ids:  []int{1, 2, 3},
			notification: models.Notification{
				EventID:  1,
				Message:  "Event Reminder",
				NotifyAt: parseTime(t, "2024-12-18 10:00:00", timeLayout),
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				// Ожидаем 3 вставки
				m.ExpectBegin()
				for _, id := range []int{1, 2, 3} {
					m.ExpectExec(`INSERT INTO NOTIFICATION`).
						WithArgs(id, 1, "Event Reminder", parseTime(t, "2024-12-18 10:00:00", timeLayout)).
						WillReturnResult(pgxmock.NewResult("INSERT", 1))
				}
				m.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name: "Ошибка при вставке одного уведомления (откат транзакции)",
			ids:  []int{1, 2, 3},
			notification: models.Notification{
				EventID:  1,
				Message:  "Event Reminder",
				NotifyAt: parseTime(t, "2024-12-18 10:00:00", timeLayout),
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				// Ожидаем начало транзакции
				m.ExpectBegin()

				// Ожидаем успешную вставку для первого и второго пользователя
				m.ExpectExec(`INSERT INTO NOTIFICATION`).
					WithArgs(1, 1, "Event Reminder", parseTime(t, "2024-12-18 10:00:00", timeLayout)).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				m.ExpectExec(`INSERT INTO NOTIFICATION`).
					WithArgs(2, 1, "Event Reminder", parseTime(t, "2024-12-18 10:00:00", timeLayout)).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				// Ошибка при вставке для третьего пользователя
				m.ExpectExec(`INSERT INTO NOTIFICATION`).
					WithArgs(3, 1, "Event Reminder", parseTime(t, "2024-12-18 10:00:00", timeLayout)).
					WillReturnError(fmt.Errorf("database error"))

				// Ожидаем откат транзакции
				m.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name: "Ошибка при начале транзакции",
			ids:  []int{1, 2},
			notification: models.Notification{
				EventID:  1,
				Message:  "Event Reminder",
				NotifyAt: parseTime(t, "2024-12-18 10:00:00", timeLayout),
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				// Симулируем ошибку при начале транзакции
				m.ExpectBegin().WillReturnError(fmt.Errorf("transaction begin error"))
			},
			expectErr: true,
		},
		{
			name: "Ошибка при коммите транзакции",
			ids:  []int{1, 2},
			notification: models.Notification{
				EventID:  1,
				Message:  "Event Reminder",
				NotifyAt: parseTime(t, "2024-12-18 10:00:00", timeLayout),
			},
			mockSetup: func(m pgxmock.PgxConnIface) {
				// Ожидаем успешную вставку для каждого пользователя
				m.ExpectBegin()
				for _, id := range []int{1, 2} {
					m.ExpectExec(`INSERT INTO NOTIFICATION`).
						WithArgs(id, 1, "Event Reminder", parseTime(t, "2024-12-18 10:00:00", timeLayout)).
						WillReturnResult(pgxmock.NewResult("INSERT", 1))
				}
				// Симулируем ошибку при коммите транзакции
				m.ExpectCommit().WillReturnError(fmt.Errorf("commit error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockConn, err := pgxmock.NewConn()
			require.NoError(t, err)
			defer mockConn.Close(ctx)

			tt.mockSetup(mockConn)

			db := repository.NewDB(mockConn)

			err = db.CreateNotificationsByUserIDs(context.Background(), tt.ids, tt.notification)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), models.LevelDB)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNotificationRepository_GetNotifications(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	timeLayout := "2006-01-02 15:04:05"

	tests := []struct {
		name         string
		userID       int
		mockSetup    func(m pgxmock.PgxConnIface)
		expectedData []models.Notification
		expectErr    bool
	}{
		{
			name:   "Успешное извлечение уведомлений",
			userID: 1,
			mockSetup: func(m pgxmock.PgxConnIface) {
				rows := pgxmock.NewRows([]string{"id", "user_id", "event_id", "notify_at", "message"}).
					AddRow(1, 1, 1, parseTime(t, "2024-12-18 10:00:00", timeLayout), "Event Reminder").
					AddRow(2, 1, 2, parseTime(t, "2024-12-19 10:00:00", timeLayout), "Another Event")
				m.ExpectQuery(`SELECT id, user_id, event_id, notify_at, message`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedData: []models.Notification{
				{ID: 1, UserID: 1, EventID: 1, NotifyAt: parseTime(t, "2024-12-18 10:00:00", timeLayout), Message: "Event Reminder"},
				{ID: 2, UserID: 1, EventID: 2, NotifyAt: parseTime(t, "2024-12-19 10:00:00", timeLayout), Message: "Another Event"},
			},
			expectErr: false,
		},
		{
			name:   "Ошибка при выполнении запроса",
			userID: 2,
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT id, user_id, event_id, notify_at, message`).
					WithArgs(2).
					WillReturnError(fmt.Errorf("query error"))
			},
			expectedData: nil,
			expectErr:    true,
		},
		{
			name:   "Ошибка при сканировании данных",
			userID: 1,
			mockSetup: func(m pgxmock.PgxConnIface) {
				// Создаем строки, но имитируем ошибку при сканировании данных
				rows := pgxmock.NewRows([]string{"id", "user_id", "event_id", "notify_at", "message"}).
					AddRow(1, 1, 1, parseTime(t, "2024-12-18 10:00:00", timeLayout), "Event Reminder")
				// Мокируем ошибку при сканировании (например, неправильный тип данных в поле)
				m.ExpectQuery(`SELECT id, user_id, event_id, notify_at, message`).
					WithArgs(1).
					WillReturnRows(rows).
					WillReturnError(fmt.Errorf("scan error"))
			},
			expectedData: nil,
			expectErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockConn, err := pgxmock.NewConn()
			require.NoError(t, err)
			defer mockConn.Close(ctx)

			tt.mockSetup(mockConn)

			db := repository.NewDB(mockConn)

			notifications, err := db.GetNotifications(ctx, tt.userID)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedData, notifications)
			}
		})
	}
}

func TestNotificationRepository_UpdateSentNotifications(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		ids       []int
		mockSetup func(m pgxmock.PgxConnIface)
		expectErr bool
	}{
		{
			name: "Успешное обновление уведомлений",
			ids:  []int{1, 2, 3},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectBegin()
				for _, id := range []int{1, 2, 3} {
					m.ExpectExec(`UPDATE notification`).
						WithArgs(id).
						WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				}
				m.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name: "Ошибка при обновлении уведомлений",
			ids:  []int{1, 2, 3},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectBegin()
				m.ExpectExec(`UPDATE notification`).
					WithArgs(1).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				m.ExpectExec(`UPDATE notification`).
					WithArgs(2).
					WillReturnError(fmt.Errorf("update error"))
				m.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name: "Ошибка при начале транзакции",
			ids:  []int{1, 2},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectBegin().WillReturnError(fmt.Errorf("transaction begin error"))
			},
			expectErr: true,
		},
		{
			name: "Ошибка при коммите транзакции",
			ids:  []int{1, 2},
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectBegin()
				m.ExpectExec(`UPDATE notification`).
					WithArgs(1).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				m.ExpectExec(`UPDATE notification`).
					WithArgs(2).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				m.ExpectCommit().WillReturnError(fmt.Errorf("commit error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockConn, err := pgxmock.NewConn()
			require.NoError(t, err)
			defer mockConn.Close(ctx)

			tt.mockSetup(mockConn)

			db := repository.NewDB(mockConn)

			err = db.UpdateSentNotifications(ctx, tt.ids)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), models.LevelDB)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func parseTime(t *testing.T, value string, layout string) time.Time {
	parsedTime, err := time.Parse(layout, value)
	require.NoError(t, err)
	return parsedTime
}
