package handlers

import (
	"testing"

	pb "kudago/internal/user/api"

	"github.com/stretchr/testify/assert"
)

func TestUserToUserResponse(t *testing.T) {
	user := &pb.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		AvatarUrl: "http://example.com/avatar.png",
	}

	userResponse := userToUserResponse(user)

	assert.Equal(t, 1, userResponse.ID)
	assert.Equal(t, "testuser", userResponse.Username)
	assert.Equal(t, "test@example.com", userResponse.Email)
	assert.Equal(t, "http://example.com/avatar.png", userResponse.ImageURL)
}

func TestWriteUsersResponse(t *testing.T) {
	users := []*pb.User{
		{ID: 1, Username: "user1", Email: "user1@example.com", AvatarUrl: "http://example.com/user1.png"},
		{ID: 2, Username: "user2", Email: "user2@example.com", AvatarUrl: "http://example.com/user2.png"},
	}

	response := writeUsersResponse(users, 2)

	assert.Len(t, response.Users, 2)
	assert.Equal(t, "user1", response.Users[0].Username)
	assert.Equal(t, "user2", response.Users[1].Username)
}
