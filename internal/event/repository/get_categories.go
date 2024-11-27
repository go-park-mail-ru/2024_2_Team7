package eventRepository

import (
	"context"
	"fmt"

	"kudago/internal/models"
)

const getCategoriesQuery = `SELECT * FROM category`

func (db *EventDB) GetCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := db.pool.Query(ctx, getCategoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", models.LevelDB, err)
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
			return nil, fmt.Errorf("%s: %w", models.LevelDB, err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}
