package repository

import (
	"testing"
	"time"

	"kudago/internal/auth/repository/auth"

	"github.com/stretchr/testify/assert"

	"kudago/internal/models"
)

func TestUserDB_nilIfEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    string
		expected *string
	}{
		{
			name:     "Empty string",
			value:    "",
			expected: nil,
		},
		{
			name:     "Non-empty string",
			value:    "http://example.com/avatar.jpg",
			expected: stringPtr("http://example.com/avatar.jpg"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := userRepository.NilIfEmpty(tt.value)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUserDB_toDomainUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		userInfo userRepository.UserInfo
		expected models.User
	}{
		{
			name: "User with Image URL",
			userInfo: userRepository.UserInfo{
				ID:        1,
				Username:  "testuser",
				Email:     "testuser@example.com",
				ImageURL:  stringPtr("http://example.com/avatar.jpg"),
				CreatedAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: models.User{
				ID:       1,
				Username: "testuser",
				Email:    "testuser@example.com",
				ImageURL: "http://example.com/avatar.jpg",
			},
		},
		{
			name: "User without Image URL",
			userInfo: userRepository.UserInfo{
				ID:        2,
				Username:  "anotheruser",
				Email:     "anotheruser@example.com",
				ImageURL:  nil,
				CreatedAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: models.User{
				ID:       2,
				Username: "anotheruser",
				Email:    "anotheruser@example.com",
				ImageURL: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := userRepository.ToDomainUser(tt.userInfo)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
