package eventRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"
)

const getSubscribersIDsQuery = `
	SELECT subscriber_id
	FROM SUBSCRIPTION
	WHERE follows_id = $1;
`

func (db EventDB) GetSubscribersIDs(ctx context.Context, ID int) ([]int, error) {
	rows, err := db.pool.Query(ctx, getSubscribersIDsQuery, ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", models.LevelDB, err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", models.LevelDB, err)
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", models.LevelDB, err)
	}

	return ids, nil
}
