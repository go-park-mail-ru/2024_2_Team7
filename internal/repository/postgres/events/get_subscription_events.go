package eventRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"
)

const getSubscriptionEventsQuery = `
	SELECT event.id, event.title, event.description, event.event_start, event.event_finish, 
		event.location, event.capacity, event.created_at, event.user_id, event.category_id, 
		COALESCE(array_agg(COALESCE(tag.name, '')), '{}') AS tags, media_url.url AS media_link
	FROM event
	INNER JOIN SUBSCRIPTION ON event.user_id = SUBSCRIPTION.subscribed_id
	LEFT JOIN event_tag ON event.id = event_tag.event_id
	LEFT JOIN tag ON tag.id = event_tag.tag_id
	LEFT JOIN media_url ON event.id = media_url.event_id
	WHERE SUBSCRIPTION.subscriber_id=$1 
	GROUP BY event.id, media_url.url
	ORDER BY event.event_finish ASC
	LIMIT $2 OFFSET $3`

func (db *EventDB) GetSubscriptionEvents(ctx context.Context, userID int, paginationParams models.PaginationParams) ([]models.Event, error) {
	rows, err := db.pool.Query(ctx, getSubscriptionEventsQuery, userID, paginationParams.Limit, paginationParams.Offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", models.LevelDB, err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var eventInfo EventInfo
		err = rows.Scan(
			&eventInfo.ID,
			&eventInfo.Title,
			&eventInfo.Description,
			&eventInfo.EventStart,
			&eventInfo.EventFinish,
			&eventInfo.Location,
			&eventInfo.Capacity,
			&eventInfo.CreatedAt,
			&eventInfo.UserID,
			&eventInfo.CategoryID,
			&eventInfo.Tags,
			&eventInfo.ImageURL,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", models.LevelDB, err)
		}

		event, err := db.toDomainEvent(ctx, eventInfo)
		if err != nil {
			continue
		}
		events = append(events, event)
	}

	return events, nil
}
