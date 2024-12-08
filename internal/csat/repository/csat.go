package repository

import (
	"context"
	"fmt"

	"kudago/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type CSATDB struct {
	pool *pgxpool.Pool
}

func NewDB(pool *pgxpool.Pool) *CSATDB {
	return &CSATDB{
		pool: pool,
	}
}

const getTestQuery = `
	SELECT 
		t.id AS test_id, 
		t.title AS test_title, 
		q.id AS question_id, 
		q.question AS question_text
	FROM test t
	LEFT JOIN question q ON t.id = q.test_id
	WHERE LOWER(t.title) ILIKE $1
	ORDER BY t.created_at DESC, q.created_at DESC;
`

func (db *CSATDB) GetTest(ctx context.Context, query string) (models.Test, error) {
	rows, err := db.pool.Query(ctx, getTestQuery, query+"%")
	if err != nil {
		return models.Test{}, errors.Wrap(err, models.LevelDB)
	}
	defer rows.Close()

	var test models.Test
	test.Questions = []models.Question{}

	for rows.Next() {
		var (
			testID       int
			testTitle    string
			questionID   *int
			questionText *string
		)

		err := rows.Scan(&testID, &testTitle, &questionID, &questionText)
		if err != nil {
			return models.Test{}, errors.Wrap(err, models.LevelDB)
		}

		if test.ID == 0 {
			test.ID = testID
			test.Title = testTitle
		}

		if questionID != nil {
			test.Questions = append(test.Questions, models.Question{
				ID:   *questionID,
				Text: *questionText,
			})
		}
	}

	if test.ID == 0 {
		return models.Test{}, models.ErrNotFound
	}

	return test, nil
}

const insertAnswersQuery = `
	INSERT INTO answers (question_id, user_id, answer) 
	VALUES ($1, $2, $3)`

func (db *CSATDB) AddAnswers(ctx context.Context, answers []models.Answer, userID int) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", models.LevelDB, err)
	}
	defer tx.Rollback(ctx)

	for _, answer := range answers {
		_, err := tx.Exec(ctx, insertAnswersQuery, answer.QuestionID, userID, answer.Value)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				return models.ErrForeignKeyViolation
			}
			return fmt.Errorf("%s: %w", models.LevelDB, err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", models.LevelDB, err)
	}
	return nil
}

const getStatsQuery = `
	SELECT 
		q.id AS question_id,
		q.question,
		AVG(a.answer) AS average_answer
	FROM   QUESTION q
	LEFT JOIN  ANSWERS a ON q.id = a.question_id
	GROUP BY   q.id, q.question
	ORDER BY  q.id;
`

func (db *CSATDB) GetStatistics(ctx context.Context) ([]models.Stats, error) {
	rows, err := db.pool.Query(ctx, getStatsQuery)
	if err != nil {
		return nil, errors.Wrap(err, models.LevelDB)
	}
	defer rows.Close()

	statistics := []models.Stats{}

	for rows.Next() {
		var (
			questionID *int
			question   *string
			value      *int
		)

		err := rows.Scan(&questionID, &question, &value)
		if err != nil {
			return nil, errors.Wrap(err, models.LevelDB)
		}

		if questionID != nil && value != nil {
			statistics = append(statistics, models.Stats{
				ID:       *questionID,
				Question: *question,
				Value:    *value,
			})
		}
	}

	return statistics, nil
}
