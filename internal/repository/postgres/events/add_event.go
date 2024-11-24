package eventRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"
)

const addEventQuery = `
	INSERT INTO event (title, description, event_start, event_finish, location, capacity, user_id, category_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id`

func (db *EventDB) AddEvent(ctx context.Context, event models.Event) (models.Event, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return models.Event{}, fmt.Errorf("%s: %w", models.LevelDB, err)
	}
	defer tx.Rollback(ctx)

	var id int
	err = tx.QueryRow(ctx, addEventQuery, event.Title, event.Description, event.EventStart, event.EventEnd, event.Location, event.Capacity, event.AuthorID, event.CategoryID).Scan(&id)
	if err != nil {
		return models.Event{}, fmt.Errorf("%s: %w", models.LevelDB, err)
	}

	event.ID = id
	err = db.addTagsToEvent(ctx, tx, id, event.Tag)
	if err != nil {
		return models.Event{}, err
	}

	err = db.addMediaURL(ctx, tx, id, event.ImageURL)
	if err != nil {
		return models.Event{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return models.Event{}, fmt.Errorf("%s: %w", models.LevelDB, err)
	}
	event.ID = id
	return event, nil
}
