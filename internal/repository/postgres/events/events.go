package eventRepository

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"kudago/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	ImageURL    *string   `db:"image"`
}

func NewDB(pool *pgxpool.Pool) *EventDB {
	return &EventDB{
		pool: pool,
	}
}

const selectUpcomingEventsQuery = `
	SELECT event.id, event.title, event.description, event.event_start, event.event_finish, 
		event.location, event.capacity, event.created_at, event.user_id, event.category_id, 
		COALESCE(array_agg(DISTINCT COALESCE(tag.name, '')), '{}') AS tags, media_url.url AS media_link
	FROM event
	LEFT JOIN event_tag ON event.id = event_tag.event_id
	LEFT JOIN tag ON tag.id = event_tag.tag_id
	LEFT JOIN media_url ON event.id = media_url.event_id
	WHERE event.event_finish >= NOW()
	GROUP BY event.id, media_url.url
	ORDER BY event.event_start ASC
	LIMIT $1 OFFSET $2`

func (db *EventDB) GetUpcomingEvents(ctx context.Context, paginationParams models.PaginationParams) ([]models.Event, error) {
	rows, err := db.pool.Query(ctx, selectUpcomingEventsQuery, paginationParams.Limit, paginationParams.Offset)
	if err != nil {
		return nil, errors.Wrap(err, models.LevelDB)
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
		)
		if err != nil {
			return nil, errors.Wrap(err, models.LevelDB)
		}
		event, err := db.toDomainEvent(ctx, eventInfo)
		if err != nil {
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

const getEventByIDQuery = `
	SELECT event.id, event.title, event.description, event.event_start, event.event_finish, 
	event.location, event.capacity, event.created_at, event.user_id, event.category_id, 
	COALESCE(array_agg(COALESCE(tag.name, '')), '{}') AS tags, media_url.url AS media_link
	FROM event
	LEFT JOIN event_tag ON event.id = event_tag.event_id
	LEFT JOIN tag ON tag.id = event_tag.tag_id
	LEFT JOIN media_url ON event.id = media_url.event_id
	WHERE event.id=$1
	GROUP BY event.id, media_url.url`

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
		&eventInfo.ImageURL,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Event{}, errors.Wrap(models.ErrEventNotFound, models.LevelDB)
		}
		return models.Event{}, errors.Wrap(err, models.LevelDB)
	}

	event, err := db.toDomainEvent(ctx, eventInfo)
	if err != nil {
		return models.Event{}, models.ErrInternal
	}
	return event, nil
}

const getEventsByCategoryQuery = `
	SELECT event.id, event.title, event.description, event.event_start, event.event_finish,
		event.location, event.capacity, event.created_at, event.user_id, event.category_id,
		COALESCE(array_agg(COALESCE(tag.name, '')), '{}') AS tags, media_url.url AS media_link
	FROM event
	LEFT JOIN event_tag ON event.id = event_tag.event_id
	LEFT JOIN tag ON tag.id = event_tag.tag_id
	LEFT JOIN media_url ON event.id = media_url.event_id
	WHERE event.category_id=$1 	AND event.event_finish >= NOW()
	GROUP BY event.id, media_url.url
	ORDER BY event.event_finish ASC
	LIMIT $2 OFFSET $3`

func (db *EventDB) GetEventsByCategory(ctx context.Context, categoryID int,paginationParams models.PaginationParams) ([]models.Event, error) {
	rows, err := db.pool.Query(ctx, getEventsByCategoryQuery, categoryID, paginationParams.Limit, paginationParams.Offset)
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
			&eventInfo.ImageURL,
		)
		if err != nil {
			return nil, errors.Wrap(err, models.LevelDB)
		}

		event, err := db.toDomainEvent(ctx, eventInfo)
		if err != nil {
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

const deleteEventQuery = `DELETE FROM event WHERE id=$1`

func (db *EventDB) DeleteEvent(ctx context.Context, ID int) error {
	_, err := db.pool.Exec(ctx, deleteEventQuery, ID)
	return errors.Wrap(err, models.LevelDB)
}

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
		return models.Event{}, errors.Wrap(err, models.LevelDB)
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
		return models.Event{}, errors.Wrap(err, models.LevelDB)
	}

	if len(updatedEvent.Tag) > 0 {
		err = db.updateTagsForEvent(ctx, tx, updatedEvent.ID, updatedEvent.Tag)
		if err != nil {
			return models.Event{}, errors.Wrap(err, models.LevelDB)
		}
	}

	if updatedEvent.ImageURL != "" {
		err = db.updateMediaURL(ctx, tx, updatedEvent.ID, updatedEvent.ImageURL)
		if err != nil {
			return models.Event{}, errors.Wrap(err, models.LevelDB)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return models.Event{}, errors.Wrap(err, models.LevelDB)
	}

	event, err := db.toDomainEvent(ctx, eventInfo)
	if err != nil {
		return models.Event{}, errors.Wrap(err, models.LevelDB)
	}
	event.Tag = updatedEvent.Tag
	event.ImageURL = updatedEvent.ImageURL
	return event, nil
}

const deleteEventTagsQuery = `DELETE FROM event_tag WHERE event_id = $1`

func (db *EventDB) updateTagsForEvent(ctx context.Context, tx pgx.Tx, eventID int, tags []string) error {
	_, err := tx.Exec(ctx, deleteEventTagsQuery, eventID)
	if err != nil {
		return errors.Wrap(err, models.LevelDB)
	}

	return db.addTagsToEvent(ctx, tx, eventID, tags)
}

const deleteMediaURLQuery = `DELETE FROM media_url WHERE event_id = $1`

func (db *EventDB) updateMediaURL(ctx context.Context, tx pgx.Tx, eventID int, imageURL string) error {
	_, err := tx.Exec(ctx, deleteMediaURLQuery, eventID)
	if err != nil {
		return errors.Wrap(err, models.LevelDB)
	}

	return db.addMediaURL(ctx, tx, eventID, imageURL)
}

const addEventQuery = `
	INSERT INTO event (title, description, event_start, event_finish, location, capacity, user_id, category_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id`

func (db *EventDB) AddEvent(ctx context.Context, event models.Event) (models.Event, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return models.Event{}, errors.Wrap(err, models.LevelDB)
	}
	defer tx.Rollback(ctx)

	var id int
	err = db.pool.QueryRow(ctx, addEventQuery, event.Title, event.Description, event.EventStart, event.EventEnd, event.Location, event.Capacity, event.AuthorID, event.CategoryID).Scan(&id)
	if err != nil {
		return models.Event{}, errors.Wrap(err, models.LevelDB)
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
		return models.Event{}, errors.Wrap(err, models.LevelDB)
	}
	return db.GetEventByID(ctx, id)
}

const insertTagsQuery = `
	INSERT INTO tag (name) 
	VALUES ($1)
	ON CONFLICT (name) DO NOTHING`

const selectTagIDsQuery = `SELECT id FROM tag WHERE name = ANY($1)`

func (db *EventDB) addTagsToEvent(ctx context.Context, tx pgx.Tx, eventID int, tags []string) error {
	tagIDs := make([]int, 0, len(tags))

	for _, tag := range tags {
		_, err := tx.Exec(ctx, insertTagsQuery, tag)
		if err != nil {
			return errors.Wrap(err, models.LevelDB)
		}
	}

	rows, err := tx.Query(ctx, selectTagIDsQuery, tags)
	if err != nil {
		return errors.Wrap(err, models.LevelDB)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return errors.Wrap(err, models.LevelDB)
		}
		tagIDs = append(tagIDs, id)
	}

	err = db.linkTagsToEvent(ctx, tx, eventID, tagIDs)
	if err != nil {
		return err
	}
	return nil
}

const insertEventTagQuery = `
	INSERT INTO event_tag (event_id, tag_id)
	VALUES ($1, $2)`

func (db *EventDB) linkTagsToEvent(ctx context.Context, tx pgx.Tx, eventID int, tagIDs []int) error {
	batch := &pgx.Batch{}
	for _, tagID := range tagIDs {
		batch.Queue(insertEventTagQuery, eventID, tagID)
	}
	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	for range tagIDs {
		_, err := br.Exec()
		if err != nil {
			return errors.Wrap(err, models.LevelDB)
		}
	}

	return nil
}

const insertMediaURL = `
	INSERT INTO media_url (event_id, url)
	VALUES ($1, $2)`

func (db *EventDB) addMediaURL(ctx context.Context, tx pgx.Tx, eventID int, imageURL string) error {
	_, err := tx.Exec(ctx, insertMediaURL, eventID, imageURL)
	if err != nil {
		return errors.Wrap(err, models.LevelDB)
	}
	return nil
}

func (db *EventDB) toDomainEvent(ctx context.Context, eventInfo EventInfo) (models.Event, error) {
	url := ""

	if eventInfo.ImageURL != nil {
		url = *eventInfo.ImageURL
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
		CategoryID:  eventInfo.CategoryID,
		ImageURL:    url,
	}, nil
}

const getCategoriesQuery = `SELECT * FROM category`

func (db *EventDB) GetCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := db.pool.Query(ctx, getCategoriesQuery)
	if err != nil {
		return nil, errors.Wrap(err, models.LevelDB)
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
			return nil, errors.Wrap(err, models.LevelDB)
		}
		categories = append(categories, category)
	}

	return categories, nil
}

const selectPastEventsQuery = `
	SELECT event.id, event.title, event.description, event.event_start, event.event_finish,
		event.location, event.capacity, event.created_at, event.user_id, event.category_id, 
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
		return nil, errors.Wrap(err, models.LevelDB)
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
		)
		if err != nil {
			return nil, errors.Wrap(err, models.LevelDB)
		}
		event, err := db.toDomainEvent(ctx, eventInfo)
		if err != nil {
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

const getEventsByUserQuery = `
	SELECT event.id, event.title, event.description, event.event_start, event.event_finish, 
		event.location, event.capacity, event.created_at, event.user_id, event.category_id, 
		COALESCE(array_agg(COALESCE(tag.name, '')), '{}') AS tags, media_url.url AS media_link
	FROM event
	LEFT JOIN event_tag ON event.id = event_tag.event_id
	LEFT JOIN tag ON tag.id = event_tag.tag_id
	LEFT JOIN media_url ON event.id = media_url.event_id
	WHERE event.user_id=$1 
	GROUP BY event.id, media_url.url
	ORDER BY event.event_finish ASC
	LIMIT $2 OFFSET $3`

func (db *EventDB) GetEventsByUser(ctx context.Context, userID int, paginationParams models.PaginationParams) ([]models.Event, error) {
	rows, err := db.pool.Query(ctx, getEventsByUserQuery, userID, paginationParams.Limit, paginationParams.Offset)
	if err != nil {
		return nil, errors.Wrap(err, models.LevelDB)
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
			return nil, errors.Wrap(err, models.LevelDB)
		}

		event, err := db.toDomainEvent(ctx, eventInfo)
		if err != nil {
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

const baseSearchQuery = `
	SELECT event.id, event.title, event.description, event.event_start, event.event_finish, 
		event.location, event.capacity, event.created_at, event.user_id, event.category_id,
		COALESCE(array_agg(DISTINCT tag.name), '{}') AS tags, media_url.url AS media_link
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
	HAVING (array_length($5::TEXT[], 1) = 0 OR array_agg(DISTINCT LOWER(tag.name)) @> $5::TEXT[])
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
	}

	rows, err := db.pool.Query(ctx, baseSearchQuery, args...)
	if err != nil {
		return nil, errors.Wrap(err, models.LevelDB)
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
		)
		if err != nil {
			return nil, errors.Wrap(err, models.LevelDB)
		}

		event, err := db.toDomainEvent(ctx, eventInfo)
		if err != nil {
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

func nilIfZero(value int) interface{} {
	if value == 0 {
		return nil
	}
	return value
}

func nilIfEmpty(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}

func tagsToArray(tags []string) interface{} {
	if len(tags) == 0 {
		return nil
	}

	return tags
}
