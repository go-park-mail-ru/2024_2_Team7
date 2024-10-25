package eventRepository

import (
	"context"
	"errors"
	"time"

	"kudago/internal/models"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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
	CategoryID  *int      `db:"category_id"`
	Tags        []string  `db:"tags"`
}

func NewDB(pool *pgxpool.Pool) *EventDB {
	return &EventDB{
		pool: pool,
	}
}

func (db *EventDB) GetAllEvents(ctx context.Context) ([]models.Event, error) {
	rawQuery := `
		SELECT event.id, event.title, event.description, event.event_start, event.event_finish, event.location, event.capacity, event.created_at, event.user_id, event.category_id, COALESCE(array_agg(COALESCE(tag.name, '')), '{}') AS tags
		FROM event
		LEFT JOIN event_tag ON event.id = event_tag.event_id
		LEFT JOIN tag ON tag.id = event_tag.tag_id
		GROUP BY event.id`

	rows, err := db.pool.Query(ctx, rawQuery)
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

func (db *EventDB) GetEventsByTag(ctx context.Context, tag string) ([]models.Event, error) {
	rawQuery := `
		SELECT event.id, event.title, event.description, event.event_start, event.event_finish, event.location, event.capacity, event.created_at, event.user_id, event.category_id, COALESCE(array_agg(COALESCE(tag.name, '')), '{}') AS tags
		FROM event
		LEFT JOIN event_tag ON event.id = event_tag.event_id
		LEFT JOIN tag ON tag.id = event_tag.tag_id
		WHERE tag.name=$1
		GROUP BY event.id`

	rows, err := db.pool.Query(ctx, rawQuery, tag)
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
	rawQuery := `
		SELECT event.id, event.title, event.description, event.event_start, event.event_finish, event.location, event.capacity, event.created_at, event.user_id, event.category_id, COALESCE(array_agg(COALESCE(tag.name, '')), '{}') AS tags
		FROM event
		LEFT JOIN event_tag ON event.id = event_tag.event_id
		LEFT JOIN tag ON tag.id = event_tag.tag_id
		WHERE event.id=$1
		GROUP BY event.id`

	var eventInfo EventInfo
	err := db.pool.QueryRow(ctx, rawQuery, ID).Scan(
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

func (db *EventDB) GetEventsByCategory(ctx context.Context, category string) ([]models.Event, error) {
	ID, err := db.getCategoryID(ctx, category)
	if err != nil {
		return nil, err
	}

	rawQuery := `
		SELECT event.id, event.title, event.description, event.event_start, event.event_finish, event.location, event.capacity, event.created_at, event.user_id, event.category_id, COALESCE(array_agg(COALESCE(tag.name, '')), '{}') AS tags
		FROM event
		LEFT JOIN event_tag ON event.id = event_tag.event_id
		LEFT JOIN tag ON tag.id = event_tag.tag_id
		WHERE event.category_id=$1
		GROUP BY event.id`

	rows, err := db.pool.Query(ctx, rawQuery, ID)
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
	query := `DELETE FROM event WHERE id=$1`

	_, err := db.pool.Exec(ctx, query, ID)
	return err
}

func (db *EventDB) UpdateEvent(ctx context.Context, updatedEvent models.Event) error {
	rawQuery := `
		UPDATE event
		SET title = $1, description = $2, event_start = $3, event_finish = $4
		WHERE id = $5`

	_, err := db.pool.Exec(ctx, rawQuery,
		updatedEvent.Title,
		updatedEvent.Description,
		updatedEvent.EventStart,
		updatedEvent.EventEnd,
		updatedEvent.ID,
	)
	return err
}

func (db *EventDB) AddEvent(ctx context.Context, event models.Event) (models.Event, error) {
	categoryID, err := db.getCategoryID(ctx, event.Category)
	if err != nil {
		return models.Event{}, err
	}

	query := `
		INSERT INTO event (title, description, event_start, event_finish, location, capacity, user_id, category_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	var id int
	err = db.pool.QueryRow(ctx, query, event.Title, event.Description, event.EventStart, event.EventEnd, event.Location, event.Capacity, event.AuthorID, categoryID).Scan(&id)
	if err != nil {
		return models.Event{}, err
	}

	event.ID = id
	err = db.AddTagsToEvent(ctx, id, event.Tag)
	if err != nil {
		return models.Event{}, err
	}
	return db.GetEventByID(ctx, id)
}

func (db *EventDB) AddTagsToEvent(ctx context.Context, eventID int, tags []string) error {
	for _, tag := range tags {
		tagID, err := db.CreateOrGetTagID(ctx, tag)
		if err != nil {
			return err
		}

		err = db.LinkTagToEvent(ctx, eventID, tagID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *EventDB) CreateOrGetTagID(ctx context.Context, tag string) (int, error) {
	var tagID int
	query := `SELECT id FROM tag WHERE name = $1`

	err := db.pool.QueryRow(ctx, query, tag).Scan(&tagID)
	if errors.Is(err, pgx.ErrNoRows) {
		insertQuery := `INSERT INTO tag (name) VALUES ($1) RETURNING id`
		err = db.pool.QueryRow(ctx, insertQuery, tag).Scan(&tagID)
	}

	if err != nil {
		return 0, err
	}
	return tagID, nil
}

func (db *EventDB) LinkTagToEvent(ctx context.Context, eventID int, tagID int) error {
	query := `INSERT INTO event_tag (event_id, tag_id) VALUES ($1, $2)`
	_, err := db.pool.Exec(ctx, query, eventID, tagID)
	return err
}

func (db *EventDB) toDomainEvent(ctx context.Context, eventInfo EventInfo) (models.Event, error) {
	var category string
	query := `SELECT name FROM category WHERE id = $1`
	err := db.pool.QueryRow(ctx, query, eventInfo.CategoryID).Scan(&category)
	if err != nil {
		return models.Event{}, err
	}

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
		Category:    category,
	}, nil
}

func (db *EventDB) getCategoryID(ctx context.Context, category string) (int, error) {
	var categoryID int
	query := `SELECT id FROM category WHERE name = $1`

	err := db.pool.QueryRow(ctx, query, category).Scan(&categoryID)
	if err == pgx.ErrNoRows {
		return 0, models.ErrInvalidCategory
	}
	return categoryID, nil
}
