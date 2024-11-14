package eventRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"
)

const insertNewFavorite = `
	INSERT INTO FAVORITE_EVENT (user_id, event_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING`

func (db *EventDB) AddEventToFavorites(ctx context.Context, newFavorite models.FavoriteEvent) error {
	_, err := db.pool.Exec(ctx, insertNewFavorite, newFavorite.UserID, newFavorite.EventID)
	if err != nil {
		return fmt.Errorf("%s: %w", models.LevelDB, err)
	}
	return nil
}
