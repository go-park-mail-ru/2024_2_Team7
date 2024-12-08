package eventRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"
)

const deleteEventQuery = `
DELETE FROM event WHERE id=$1`

func (db *EventDB) DeleteEvent(ctx context.Context, ID int) error {
	_, err := db.pool.Exec(ctx, deleteEventQuery, ID)
	if err != nil {
		return fmt.Errorf("%s: %w", models.LevelDB, err)
	}
	return nil
}
