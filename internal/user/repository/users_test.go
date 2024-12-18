package userRepository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"kudago/internal/models"
)

func TestNilIfEmpty(t *testing.T) {
	// Тестируем функцию NilIfEmpty
	tests := []struct {
		input    string
		expected *string
	}{
		{"", nil},
		{"non-empty", stringPtr("non-empty")},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actual := NilIfEmpty(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestToDomainUser(t *testing.T) {
	// Тестируем функцию ToDomainUser
	tests := []struct {
		input    UserInfo
		expected models.User
	}{
		{
			input: UserInfo{
				ID:         1,
				Username:   "user1",
				Email:      "user1@example.com",
				ImageURL:   stringPtr("http://example.com/avatar1"),
				CreatedAt:  time.Now(),
				ModifiedAt: time.Now(),
			},
			expected: models.User{
				ID:       1,
				Username: "user1",
				Email:    "user1@example.com",
				ImageURL: "http://example.com/avatar1",
			},
		},
		{
			input: UserInfo{
				ID:         2,
				Username:   "user2",
				Email:      "user2@example.com",
				ImageURL:   nil,
				CreatedAt:  time.Now(),
				ModifiedAt: time.Now(),
			},
			expected: models.User{
				ID:       2,
				Username: "user2",
				Email:    "user2@example.com",
				ImageURL: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input.Username, func(t *testing.T) {
			actual := ToDomainUser(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
