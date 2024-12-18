package eventRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"
)

const selectUserIDsByFavoriteEvent = `
	SELECT user_id FROM FAVORITE_EVENT 
	WHERE event_id=$1`

func (db *EventDB) GetUserIDsByFavoriteEvent(ctx context.Context, eventID int) ([]int, error) {
	rows, err := db.pool.Query(ctx, selectUserIDsByFavoriteEvent, eventID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", models.LevelDB, err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var temp int
		err = rows.Scan(
			&temp,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", models.LevelDB, err)
		}

		ids = append(ids, temp)
	}

	return ids, nil
}
