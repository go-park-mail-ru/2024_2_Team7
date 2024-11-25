package eventRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"
)

const selectPastEventsQuery = `
	SELECT event.id, event.title, event.description, event.event_start, event.event_finish,
		event.location, event.capacity, event.created_at, event.user_id, event.category_id, event.lat, event.lon, 
		COALESCE(array_agg(DISTINCT COALESCE(tag.name, '')), '{}') AS tags, media_url.url AS media_link
	FROM event
	LEFT JOIN event_tag ON event.id = event_tag.event_id
	LEFT JOIN tag ON tag.id = event_tag.tag_id
	LEFT JOIN media_url ON event.id = media_url.event_id
	WHERE event.event_finish < NOW()
	GROUP BY event.id, media_url.url
	ORDER BY event.event_start DESC
	LIMIT $1 OFFSET $2`

func (db *EventDB) GetPastEvents(ctx context.Context, paginationParams models.PaginationParams) ([]models.Event, error) {
	rows, err := db.pool.Query(ctx, selectPastEventsQuery, paginationParams.Limit, paginationParams.Offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", models.LevelDB, err)
	}
	defer rows.Close()

	events := make([]models.Event, 0, paginationParams.Limit)
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
			&eventInfo.Latitude,
			&eventInfo.Longitude,
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
