package eventRepository

import (
	"context"
	"fmt"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kudago/internal/models"
)

func TestGetCategories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		mockSetup func(m pgxmock.PgxConnIface)
		expectErr bool
		expected  []models.Category
	}{
		{
			name: "успешное выполнение",
			mockSetup: func(m pgxmock.PgxConnIface) {
				rows := m.NewRows([]string{"id", "name"}).
					AddRow(1, "Music").
					AddRow(2, "Theater").
					AddRow(3, "Cinema")
				m.ExpectQuery(`SELECT \* FROM category`).WillReturnRows(rows)
			},
			expectErr: false,
			expected: []models.Category{
				{ID: 1, Name: "Music"},
				{ID: 2, Name: "Theater"},
				{ID: 3, Name: "Cinema"},
			},
		},
		{
			name: "нет категорий",
			mockSetup: func(m pgxmock.PgxConnIface) {
				rows := m.NewRows([]string{"id", "name"})
				m.ExpectQuery(`SELECT \* FROM category`).WillReturnRows(rows)
			},
			expectErr: false,
			expected:  []models.Category{},
		},
		{
			name: "ошибка запроса",
			mockSetup: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(`SELECT \* FROM category`).WillReturnError(fmt.Errorf("database error"))
			},
			expectErr: true,
			expected:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockConn, err := pgxmock.NewConn()
			require.NoError(t, err)
			defer mockConn.Close(context.Background())

			tt.mockSetup(mockConn)

			db := &EventDB{pool: mockConn}

			categories, err := db.GetCategories(context.Background())

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), models.LevelDB)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, categories)
			}
		})
	}
}
