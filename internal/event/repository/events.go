//go:generate mockgen -source /home/ksu/go/pkg/mod/github.com/jackc/pgx/v5@v5.7.1/tx.go -destination=./mocks/mocks_tx.go -package=mocks
//go:generate mockgen -source events.go -destination=./mocks/mocks.go -package=mocks

package eventRepository

import (
	"context"
	"fmt"
	"time"

	"kudago/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type EventDB struct {
	pool Pool
}

type Pool interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
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
	Latitude    float64   `db:"lat"`
	Longitude   float64   `db:"lon"`
	Tags        []string  `db:"tags"`
	ImageURL    *string   `db:"image"`
}

func NewDB(pool Pool) *EventDB {
	return &EventDB{
		pool: pool,
	}
}

const deleteEventTagsQuery = `DELETE FROM event_tag WHERE event_id = $1`

func (db *EventDB) updateTagsForEvent(ctx context.Context, tx pgx.Tx, eventID int, tags []string) error {
	_, err := tx.Exec(ctx, deleteEventTagsQuery, eventID)
	if err != nil {
		return fmt.Errorf("%s: %w", models.LevelDB, err)
	}

	return db.addTagsToEvent(ctx, tx, eventID, tags)
}

const deleteMediaURLQuery = `DELETE FROM media_url WHERE event_id = $1`

func (db *EventDB) updateMediaURL(ctx context.Context, tx pgx.Tx, eventID int, imageURL string) error {
	_, err := tx.Exec(ctx, deleteMediaURLQuery, eventID)
	if err != nil {
		return fmt.Errorf("%s: %w", models.LevelDB, err)
	}

	return db.addMediaURL(ctx, tx, eventID, imageURL)
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
			return fmt.Errorf("%s: %w", models.LevelDB, err)
		}
	}

	rows, err := tx.Query(ctx, selectTagIDsQuery, tags)
	if err != nil {
		return fmt.Errorf("%s: %w", models.LevelDB, err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("%s: %w", models.LevelDB, err)
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
			return fmt.Errorf("%s: %w", models.LevelDB, err)
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
		return fmt.Errorf("%s: %w", models.LevelDB, err)
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
		Longitude:   eventInfo.Longitude,
		Latitude:    eventInfo.Latitude,
	}, nil
}

func nilIfZero(value int) interface{} {
	if value == 0 {
		return nil
	}
	return value
}

func nilIfZeroFloat(value float64) interface{} {
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
