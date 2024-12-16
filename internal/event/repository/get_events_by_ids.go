package eventRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"
)

const getEventsByIDsQuery = `
	SELECT event.id, event.title, event.description, event.event_start, event.event_finish, 
	event.location, event.capacity, event.created_at, event.user_id, event.category_id, event.lat, event.lon, 
	COALESCE(array_agg(COALESCE(tag.name, '')), '{}') AS tags, media_url.url AS media_link
	FROM event
	LEFT JOIN event_tag ON event.id = event_tag.event_id
	LEFT JOIN tag ON tag.id = event_tag.tag_id
	LEFT JOIN media_url ON event.id = media_url.event_id
	WHERE event.id=ANY($1)
	GROUP BY event.id, media_url.url`

func (db *EventDB) GetEventsByIDs(ctx context.Context, IDs []int) ([]models.Event, error) {
	rows, err := db.pool.Query(ctx, getEventsByIDsQuery, IDs)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", models.LevelDB, err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var eventInfo EventInfo
		err := rows.Scan(
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
			&eventInfo.Latitude,
			&eventInfo.Longitude,
			&eventInfo.Tags,
			&eventInfo.ImageURL,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", models.LevelDB, err)
		}

		event, err := db.toDomainEvent(ctx, eventInfo)
		if err != nil {
			return nil, models.ErrInternal
		}
		events = append(events, event)
	}
	return events, nil
}