package eventRepository

import (
	"context"
	"fmt"
	"time"

	"kudago/internal/models"
)

const updateEventQuery = `
	UPDATE event
	SET 
		title = COALESCE($2, title), 
		description = COALESCE($3, description), 
		event_start = COALESCE($4, event_start), 
		event_finish = COALESCE($5, event_finish), 
		location = COALESCE($6, location), 
		capacity = COALESCE($7, capacity), 
		category_id = COALESCE($8, category_id), 
		updated_at = $9
	WHERE id = $1
	RETURNING id, title, description, event_start, event_finish, location, capacity, category_id, user_id
`

func (db *EventDB) UpdateEvent(ctx context.Context, updatedEvent models.Event) (models.Event, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return models.Event{}, fmt.Errorf("%s: %w", models.LevelDB, err)
	}
	defer tx.Rollback(ctx)

	var eventInfo EventInfo
	err = tx.QueryRow(ctx, updateEventQuery,
		updatedEvent.ID,
		nilIfEmpty(updatedEvent.Title),
		nilIfEmpty(updatedEvent.Description),
		nilIfEmpty(updatedEvent.EventStart),
		nilIfEmpty(updatedEvent.EventEnd),
		nilIfEmpty(updatedEvent.Location),
		nilIfZero(updatedEvent.Capacity),
		nilIfZero(updatedEvent.CategoryID),
		time.Now(),
	).Scan(
		&eventInfo.ID,
		&eventInfo.Title,
		&eventInfo.Description,
		&eventInfo.EventStart,
		&eventInfo.EventFinish,
		&eventInfo.Location,
		&eventInfo.Capacity,
		&eventInfo.CategoryID,
		&eventInfo.UserID,
	)
	
	if err != nil {
		return models.Event{}, fmt.Errorf("%s: %w", models.LevelDB, err)
	}

	if len(updatedEvent.Tag) > 0 {
		err = db.updateTagsForEvent(ctx, tx, updatedEvent.ID, updatedEvent.Tag)
		if err != nil {
			return models.Event{}, fmt.Errorf("%s: %w", models.LevelDB, err)
		}
	}

	if updatedEvent.ImageURL != "" {
		err = db.updateMediaURL(ctx, tx, updatedEvent.ID, updatedEvent.ImageURL)
		if err != nil {
			return models.Event{}, fmt.Errorf("%s: %w", models.LevelDB, err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return models.Event{}, fmt.Errorf("%s: %w", models.LevelDB, err)
	}

	event, err := db.toDomainEvent(ctx, eventInfo)
	if err != nil {
		return models.Event{}, fmt.Errorf("%s: %w", models.LevelDB, err)
	}

	event.Tag = updatedEvent.Tag
	event.ImageURL = updatedEvent.ImageURL

	return event, nil
}
