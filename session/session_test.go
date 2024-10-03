package session

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name   string
	setup  func(*SessionDB, *http.Request)
	output bool
}

func TestCheckSession(t *testing.T) {
	t.Parallel()
	db := NewSessionDB()
	username := "test_user"
	validSession := db.CreateSession(username)

	req, _ := http.NewRequest("GET", "/", nil)

	testCases := []testCase{
		{
			name: "Valid input",
			setup: func(db *SessionDB, req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  SessionToken,
					Value: validSession.Token,
				})
			},
			output: true,
		},
		{
			name: "Expired session",
			setup: func(db *SessionDB, req *http.Request) {
				validSession.Expires = time.Now().Add(-time.Minute)
				db.sessions[validSession.Token] = *validSession

				req.AddCookie(&http.Cookie{
					Name:  SessionToken,
					Value: validSession.Token,
				})
			},
			output: false,
		},
		{
			name: "No session token",
			setup: func(db *SessionDB, req *http.Request) {
			},
			output: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(db, req)

			session, exists := db.CheckSession(req)
			assert.Equal(t, tc.output, exists)

			if exists {
				assert.Equal(t, username, session.Username)
				assert.WithinDuration(t, validSession.Expires, session.Expires, time.Second)
			}
		})
	}
}
