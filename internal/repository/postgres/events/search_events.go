package eventRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"
)

/*
const baseSearchQuery = `
SELECT event.id, event.title, event.description, event.event_start, event.event_finish,

	event.location, event.capacity, event.created_at, event.user_id, event.category_id,

COALESCE(array_agg(DISTINCT tag.name) FILTER (WHERE tag.name IS NOT NULL), ARRAY[]::TEXT[]) AS tags,

	COALESCE(media_url.url, '') AS media_link

FROM event
LEFT JOIN event_tag ON event.id = event_tag.event_id
LEFT JOIN tag ON tag.id = event_tag.tag_id
LEFT JOIN media_url ON event.id = media_url.event_id
WHERE

	($1::TEXT IS NULL OR event.title ILIKE '%' || $1 || '%' OR event.description ILIKE '%' || $1 || '%')
	AND ($2::INT IS NULL OR event.category_id = $2)
	AND ($3::TIMESTAMP IS NULL OR event.event_start >= $3)
	AND ($4::TIMESTAMP IS NULL OR event.event_finish <= $4)

GROUP BY event.id, media_url.url
HAVING (

	$5::TEXT[] IS NULL
	OR array_length($5::TEXT[], 1) = 0
	OR array_length(array_agg(DISTINCT LOWER(tag.name)), 1) = 0
	OR array_agg(DISTINCT LOWER(tag.name)) @> $5::TEXT[]

)
ORDER BY event.event_finish ASC
LIMIT $6 OFFSET $7;
`
*/
const baseSearchQuery = `
    SELECT event.id, event.title, event.description, event.event_start, event.event_finish, 
           event.location, event.capacity, event.created_at, event.user_id, event.category_id, event.lat, event.lon,
           COALESCE(array_agg(DISTINCT tag.name) FILTER (WHERE tag.name IS NOT NULL), ARRAY[]::TEXT[]) AS tags,
           COALESCE(media_url.url, '') AS media_link	
    FROM event
    LEFT JOIN event_tag ON event.id = event_tag.event_id
    LEFT JOIN tag ON tag.id = event_tag.tag_id
    LEFT JOIN media_url ON event.id = media_url.event_id
    WHERE
        ($1::TEXT IS NULL OR event.title ILIKE '%' || $1 || '%' OR event.description ILIKE '%' || $1 || '%')
        AND ($2::INT IS NULL OR event.category_id = $2)
        AND ($3::TIMESTAMP IS NULL OR event.event_start >= $3)
        AND ($4::TIMESTAMP IS NULL OR event.event_finish <= $4)
        AND ($8::DOUBLE PRECISION IS NULL OR event.lat >= $8) -- Минимальная широта
        AND ($9::DOUBLE PRECISION IS NULL OR event.lat <= $9) -- Максимальная широта
        AND ($10::DOUBLE PRECISION IS NULL OR event.lon >= $10) -- Минимальная долгота
        AND ($11::DOUBLE PRECISION IS NULL OR event.lon <= $11) -- Максимальная долгота
    GROUP BY event.id, media_url.url
    HAVING (
        $5::TEXT[] IS NULL 
        OR array_length($5::TEXT[], 1) = 0 
        OR array_length(array_agg(DISTINCT LOWER(tag.name)), 1) = 0 
        OR array_agg(DISTINCT LOWER(tag.name)) @> $5::TEXT[]
    )
    ORDER BY event.event_finish ASC
    LIMIT $6 OFFSET $7;
`

func (db *EventDB) SearchEvents(ctx context.Context, params models.SearchParams, paginationParams models.PaginationParams) ([]models.Event, error) {
	args := []interface{}{
		nilIfEmpty(params.Query),
		nilIfZero(params.Category),
		nilIfEmpty(params.EventStart),
		nilIfEmpty(params.EventEnd),
		tagsToArray(params.Tags),
		paginationParams.Limit,
		paginationParams.Offset,
		nilIfZeroFloat(params.LatitudeMin),
		nilIfZeroFloat(params.LatitudeMax),
		nilIfZeroFloat(params.LongitudeMin),
		nilIfZeroFloat(params.LongitudeMax),
	}

	rows, err := db.pool.Query(ctx, baseSearchQuery, args...)
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
