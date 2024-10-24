package eventRepository

import (
	"context"
	"errors"
	"time"

	"kudago/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	eventColumns = []string{"id", "title", "description", "event_start", "event_finish", "location", "capacity", "created_at", "user_id", "category_id"}
	eventTable   = `event`
	// tagColumns   = []string{"id", "name"}
	tagTable      = `tag`
	categoryTable = `category`
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

func (db EventDB) GetAllEvents(ctx context.Context) ([]models.Event, error) {
	query := sq.Select("event.id", "event.title", "event.description", "event.event_start", "event.event_finish", "event.location", "event.capacity", "event.created_at", "event.user_id", "event.category_id", "COALESCE(array_agg(COALESCE(tag.name, '')), '{}') AS tags").
		From(eventTable).
		LeftJoin("event_tag ON event.id = event_tag.event_id").
		LeftJoin("tag ON tag.id = event_tag.tag_id").
		GroupBy("event.id").
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := db.pool.Query(ctx, rawQuery, args...)
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

func (db EventDB) GetEventsByTag(ctx context.Context, tag string) ([]models.Event, error) {
	query := sq.Select("event.id", "event.title", "event.description", "event.event_start", "event.event_finish", "event.location", "event.capacity", "event.created_at", "event.user_id", "event.category_id", "COALESCE(array_agg(COALESCE(tag.name, '')), '{}') AS tags").
		From("event").
		Join("event_tag ON event.id = event_tag.event_id").
		Join("tag ON tag.id = event_tag.tag_id").
		Where(sq.Eq{"tag.name": tag}).
		GroupBy("event.id").
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := db.pool.Query(ctx, rawQuery, args...)
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

func (db EventDB) GetEventByID(ctx context.Context, ID int) (models.Event, error) {
	query := sq.Select("event.id", "event.title", "event.description", "event.event_start", "event.event_finish", "event.location", "event.capacity", "event.created_at", "event.user_id", "event.category_id", "COALESCE(array_agg(COALESCE(tag.name, '')), '{}') AS tags").
		From(eventTable).
		LeftJoin("event_tag ON event.id = event_tag.event_id").
		LeftJoin("tag ON tag.id = event_tag.tag_id").
		Where(sq.Eq{"event.id": ID}).
		GroupBy("event.id").
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return models.Event{}, err
	}

	var eventInfo EventInfo
	err = db.pool.QueryRow(ctx, rawQuery, args...).Scan(
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

func (db EventDB) GetEventsByCategory(ctx context.Context, category string) ([]models.Event, error) {
	ID, err := db.getCategoryID(ctx, category)
	if err != nil {
		return nil, err
	}

	query := sq.Select("event.id", "event.title", "event.description", "event.event_start", "event.event_finish", "event.location", "event.capacity", "event.created_at", "event.user_id", "event.category_id", "COALESCE(array_agg(COALESCE(tag.name, '')), '{}') AS tags").
		From(eventTable).
		LeftJoin("event_tag ON event.id = event_tag.event_id").
		LeftJoin("tag ON tag.id = event_tag.tag_id").
		Where(sq.Eq{"event.category_id": ID}).
		GroupBy("event.id").
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := db.pool.Query(ctx, rawQuery, args...)
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
	query := sq.Delete(eventTable).
		Where(sq.Eq{"id": ID}).
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.pool.Exec(ctx, rawQuery, args...)
	return err
}

func (db *EventDB) UpdateEvent(ctx context.Context, updatedEvent models.Event) error {
	query := sq.Update(eventTable).
		Set("title", updatedEvent.Title).
		Set("description", updatedEvent.Description).
		Set("event_start", updatedEvent.EventStart).
		Set("event_finish", updatedEvent.EventEnd).
		Where(sq.Eq{"id": updatedEvent.ID}).
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.pool.Exec(ctx, rawQuery, args...)
	return err
}

func (db *EventDB) AddEvent(ctx context.Context, event models.Event) (models.Event, error) {
	categoryID, err := db.getCategoryID(ctx, event.Category)
	if err != nil {
		return models.Event{}, err
	}

	insertQuery := sq.Insert(eventTable).
		Columns("title", "description", "event_start", "event_finish", "location", "capacity", "user_id", "category_id").
		Values(event.Title, event.Description, event.EventStart, event.EventEnd, event.Location, event.Capacity, event.AuthorID, categoryID).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := insertQuery.ToSql()
	if err != nil {
		return models.Event{}, err
	}

	var id int
	err = db.pool.QueryRow(ctx, rawQuery, args...).Scan(&id)
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
	query := sq.Select("id").
		From(tagTable).
		Where(sq.Eq{"name": tag}).
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	err = db.pool.QueryRow(ctx, rawQuery, args...).Scan(&tagID)
	if errors.Is(err, pgx.ErrNoRows) {
		insertQuery := sq.Insert(tagTable).
			Columns("name").
			Values(tag).
			Suffix("RETURNING id").
			PlaceholderFormat(sq.Dollar)

		rawQuery, args, err = insertQuery.ToSql()
		if err != nil {
			return 0, err
		}

		err = db.pool.QueryRow(ctx, rawQuery, args...).Scan(&tagID)
		if err != nil {
			return 0, err
		}
	}

	if err != nil {
		return 0, err
	}
	return tagID, nil
}

func (db *EventDB) LinkTagToEvent(ctx context.Context, eventID int, tagID int) error {
	query := sq.Insert("event_tag").
		Columns("event_id", "tag_id").
		Values(eventID, tagID).
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.pool.Exec(ctx, rawQuery, args...)
	return err
}

func (db *EventDB) toDomainEvent(ctx context.Context, eventInfo EventInfo) (models.Event, error) {
	var category string
	query := sq.Select("name").
		From(categoryTable).
		Where(sq.Eq{"id": eventInfo.CategoryID}).
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return models.Event{}, err
	}

	err = db.pool.QueryRow(ctx, rawQuery, args...).Scan(&category)
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
	query := sq.Select("id").
		From(categoryTable).
		Where(sq.Eq{"name": category}).
		PlaceholderFormat(sq.Dollar)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	err = db.pool.QueryRow(ctx, rawQuery, args...).Scan(&categoryID)
	if err == pgx.ErrNoRows {
		return 0, models.ErrInvalidCategory
	}
	return categoryID, nil
}
