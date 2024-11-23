package repository

import (
	"context"

	"kudago/internal/models"

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
