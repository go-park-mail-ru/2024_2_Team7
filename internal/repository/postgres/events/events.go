package eventRepository

import (
	"context"
	"errors"
	"time"

	"kudago/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	selectAllEventsQuery = `
	SELECT event.id, event.title, event.description, event.event_start, event.event_finish, event.location, event.capacity, event.created_at, event.user_id, event.category_id, COALESCE(array_agg(COALESCE(tag.name, '')), '{}') AS tags
	FROM event
	LEFT JOIN event_tag ON event.id = event_tag.event_id
	LEFT JOIN tag ON tag.id = event_tag.tag_id
	GROUP BY event.id
	LIMIT $1 OFFSET $2`

	getEventsByTagsQuery = `
	WITH matching_events AS (
		SELECT event.id
		FROM event
		JOIN event_tag ON event.id = event_tag.event_id
		JOIN tag ON tag.id = event_tag.tag_id
		WHERE tag.name = ANY($1)
		GROUP BY event.id
		HAVING COUNT(DISTINCT CASE WHEN tag.name = ANY($1) THEN tag.name END) = $2
	)
	SELECT event.id, event.title, event.description, event.event_start, event.event_finish, 
		   event.location, event.capacity, event.created_at, event.user_id, event.category_id,
		   COALESCE(array_agg(DISTINCT tag.name), '{}') AS tags
	FROM event
	JOIN matching_events ON event.id = matching_events.id
	JOIN event_tag ON event.id = event_tag.event_id
	JOIN tag ON tag.id = event_tag.tag_id
	GROUP BY event.id`

	getEventByIDQuery = `
	SELECT event.id, event.title, event.description, event.event_start, event.event_finish, event.location, event.capacity, event.created_at, event.user_id, event.category_id, COALESCE(array_agg(COALESCE(tag.name, '')), '{}') AS tags
	FROM event
	LEFT JOIN event_tag ON event.id = event_tag.event_id
	LEFT JOIN tag ON tag.id = event_tag.tag_id
	WHERE event.id=$1
	GROUP BY event.id`

	getEventsByCategoryQuery = `
	SELECT event.id, event.title, event.description, event.event_start, event.event_finish, event.location, event.capacity, event.created_at, event.user_id, event.category_id, COALESCE(array_agg(COALESCE(tag.name, '')), '{}') AS tags
	FROM event
	LEFT JOIN event_tag ON event.id = event_tag.event_id
	LEFT JOIN tag ON tag.id = event_tag.tag_id
	WHERE event.category_id=$1
	GROUP BY event.id`

	updateEventQuery = `
	UPDATE event
	SET title = $1, description = $2, event_start = $3, event_finish = $4, updated_at=$5
	WHERE id = $6`

	addEventQuery = `
	INSERT INTO event (title, description, event_start, event_finish, location, capacity, user_id, category_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id`

	deleteEventQuery = `DELETE FROM event WHERE id=$1`

	insertEventTagQuery = `
	INSERT INTO event_tag (event_id, tag_id)
	VALUES ($1, $2)`

	insertTagsQuery = `
	INSERT INTO tag (name) 
	VALUES ($1)
	ON CONFLICT (name) DO NOTHING`

	getCategoriesQuery = `SELECT * FROM category`

	selectTagIDsQuery = `SELECT id FROM tag WHERE name = ANY($1)`
)

type EventDB struct {
	pool *pgxpool.Pool
}

type EventInfo struct {
	ID          int       `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	EventStart  time.Time `db:"event_start"`
	EventFinish time.Time `db:"event_finish"`
	Location    string    `db:"location"`
	Capacity    int       `db:"capacity"`
	CreatedAt   time.Time `db:"created_at"`
	UserID      int       `db:"user_id"`
	CategoryID  int       `db:"category_id"`
	Tags        []string  `db:"tags"`
}

func NewDB(pool *pgxpool.Pool) *EventDB {
	return &EventDB{
		pool: pool,
	}
}

func (db *EventDB) GetAllEvents(ctx context.Context, offset, limit int) ([]models.Event, error) {
	rows, err := db.pool.Query(ctx, selectAllEventsQuery, limit, offset)
	if err != nil {
		return nil, err
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
		)
		if err != nil {
			return nil, err
		}
		event, err := db.toDomainEvent(ctx, eventInfo)
		if err != nil {
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

func (db *EventDB) GetEventsByTags(ctx context.Context, tags []string) ([]models.Event, error) {
	rows, err := db.pool.Query(ctx, getEventsByTagsQuery, tags, len(tags))
	if err != nil {
		return nil, err
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
		)
		if err != nil {
			return nil, err
		}

		event, err := db.toDomainEvent(ctx, eventInfo)
		if err != nil {
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

func (db *EventDB) GetEventByID(ctx context.Context, ID int) (models.Event, error) {
	var eventInfo EventInfo
	err := db.pool.QueryRow(ctx, getEventByIDQuery, ID).Scan(
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
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Event{}, models.ErrEventNotFound
		}
		return models.Event{}, err
	}

	event, err := db.toDomainEvent(ctx, eventInfo)
	if err != nil {
		return models.Event{}, models.ErrInternal
	}
	return event, nil
}

func (db *EventDB) GetEventsByCategory(ctx context.Context, categoryID int) ([]models.Event, error) {
	rows, err := db.pool.Query(ctx, getEventsByCategoryQuery, categoryID)
	if err != nil {
		return nil, err
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
		)
		if err != nil {
			return nil, err
		}

		event, err := db.toDomainEvent(ctx, eventInfo)
		if err != nil {
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

func (db *EventDB) DeleteEvent(ctx context.Context, ID int) error {
	_, err := db.pool.Exec(ctx, deleteEventQuery, ID)
	return err
}

func (db *EventDB) UpdateEvent(ctx context.Context, updatedEvent models.Event) error {
	_, err := db.pool.Exec(ctx, updateEventQuery,
		updatedEvent.Title,
		updatedEvent.Description,
		updatedEvent.EventStart,
		updatedEvent.EventEnd,
		time.Now(),
		updatedEvent.ID,
	)
	return err
}

func (db *EventDB) AddEvent(ctx context.Context, event models.Event) (models.Event, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return models.Event{}, err
	}
	defer tx.Rollback(ctx)

	var id int
	err = db.pool.QueryRow(ctx, addEventQuery, event.Title, event.Description, event.EventStart, event.EventEnd, event.Location, event.Capacity, event.AuthorID, event.CategoryID).Scan(&id)
	if err != nil {
		return models.Event{}, err
	}

	event.ID = id
	err = db.AddTagsToEvent(ctx, tx, id, event.Tag)
	if err != nil {
		return models.Event{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return models.Event{}, err
	}
	return db.GetEventByID(ctx, id)
}

func (db *EventDB) AddTagsToEvent(ctx context.Context, tx pgx.Tx, eventID int, tags []string) error {
	tagIDs := make([]int, 0, len(tags))

	for _, tag := range tags {
		_, err := tx.Exec(ctx, insertTagsQuery, tag)
		if err != nil {
			return err
		}
	}

	rows, err := tx.Query(ctx, selectTagIDsQuery, tags)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		tagIDs = append(tagIDs, id)
	}

	err = db.LinkTagsToEvent(ctx, tx, eventID, tagIDs)
	if err != nil {
		return err
	}
	return nil
}

func (db *EventDB) LinkTagsToEvent(ctx context.Context, tx pgx.Tx, eventID int, tagIDs []int) error {
	batch := &pgx.Batch{}
	for _, tagID := range tagIDs {
		batch.Queue(insertEventTagQuery, eventID, tagID)
	}
	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	for range tagIDs {
		_, err := br.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *EventDB) toDomainEvent(ctx context.Context, eventInfo EventInfo) (models.Event, error) {
	return models.Event{
		ID:          eventInfo.ID,
		Title:       eventInfo.Title,
		Description: eventInfo.Description,
		EventStart:  eventInfo.EventStart.Format(time.RFC3339),
		EventEnd:    eventInfo.EventFinish.Format(time.RFC3339),
		AuthorID:    eventInfo.UserID,
		Tag:         eventInfo.Tags,
		Location:    eventInfo.Location,
		Capacity:    eventInfo.Capacity,
		CategoryID:  eventInfo.CategoryID,
	}, nil
}

func (db *EventDB) GetCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := db.pool.Query(ctx, getCategoriesQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0, 10)
	for rows.Next() {
		var category models.Category
		err = rows.Scan(
			&category.ID,
			&category.Name,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}
