package repository

import (
	"context"
	"fmt"

	"kudago/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationDB struct {
	pool *pgxpool.Pool
}

func NewDB(pool *pgxpool.Pool) *NotificationDB {
	return &NotificationDB{
		pool: pool,
	}
}

const getNotificationsQuery = `
        SELECT id, user_id, event_id, notify_at, message
        FROM notification
        WHERE notify_at <= NOW() AND is_sent = FALSE AND user_id=$1
    `

func (s *NotificationDB) GetNotifications(ctx context.Context, userID int) ([]models.Notification, error) {
	rows, err := s.pool.Query(ctx, getNotificationsQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.EventID, &n.NotifyAt, &n.Message); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

const updateSentNotificationsQuery = `
    UPDATE notification
    SET is_sent = TRUE
    WHERE id = $1
`

func (db *NotificationDB) UpdateSentNotifications(ctx context.Context, ids []int) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", models.LevelDB, err)
	}
	defer tx.Rollback(ctx)

	for _, id := range ids {
		_, err = db.pool.Exec(ctx, updateSentNotificationsQuery, id)
		if err != nil {
			return fmt.Errorf("%s: %w", models.LevelDB, err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", models.LevelDB, err)
	}

	return nil
}

const deleteNotificationQuery = `DELETE FROM NOTIFICATION WHERE id=$1`

func (db *NotificationDB) DeleteNotification(ctx context.Context, ID int) error {
	_, err := db.pool.Exec(ctx, deleteNotificationQuery, ID)
	if err != nil {
		return fmt.Errorf("%s: %w", models.LevelDB, err)
	}
	return nil
}

const createNotificationQuery = `
	INSERT INTO NOTIFICATION (user_id, event_id, message, notify_at)
	VALUES ($1, $2, $3, $4)
	`

func (db *NotificationDB) CreateNotification(ctx context.Context, notification models.Notification) error {
	_, err := db.pool.Exec(ctx, createNotificationQuery,
		notification.UserID,
		notification.EventID,
		notification.Message,
		notification.NotifyAt,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", models.LevelDB, err)
	}

	return nil
}

const createNotificationsByUserIDsQuery = `
	INSERT INTO NOTIFICATION (user_id, event_id, message, notify_at)
	VALUES ($1, $2, $3, $4)
	`

func (db *NotificationDB) CreateNotificationsByUserIDs(ctx context.Context, ids []int, ntf models.Notification) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", models.LevelDB, err)
	}
	defer tx.Rollback(ctx)

	for _, id := range ids {
		_, err = db.pool.Exec(ctx, createNotificationsByUserIDsQuery, id, ntf.EventID, ntf.Message, ntf.NotifyAt)
		if err != nil {
			return fmt.Errorf("%s: %w", models.LevelDB, err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", models.LevelDB, err)
	}

	return nil
}
