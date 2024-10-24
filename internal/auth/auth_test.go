package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"kudago/internal/users"
	"kudago/session"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name             string
	req              *http.Request
	setup            func()
	expectedUsername string
	expectedStatus   int
}

func setupTest() *Handler {
	handler := &Handler{
		UserDB:    *users.NewUserDB(),
		SessionDb: *session.NewSessionDB(),
	}
	return handler
}

func TestRegister(t *testing.T) {
	handler := setupTest()
	newUser := users.User{
		Username: "new_user",
		Email:    "new_user@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(newUser)

	testCases := []testCase{
		{
			name: "Valid Input",
			req:  httptest.NewRequest("POST", "/register", bytes.NewReader(body)),
			setup: func() {
			},
			expectedStatus:   http.StatusCreated,
			expectedUsername: newUser.Username,
		},
		{
			name: "User already exists",
			req:  httptest.NewRequest("POST", "/register", bytes.NewReader(body)),
			setup: func() {
				handler.UserDB.AddUser(&newUser)
			},
			expectedStatus:   http.StatusConflict,
			expectedUsername: "",
		},
		{
			name: "Email already used",
			req:  httptest.NewRequest("POST", "/register", bytes.NewReader(body)),
			setup: func() {
				handler.UserDB.AddUser(&users.User{
					Username: "new_user2",
					Email:    "new_user@example.com",
					Password: "password123",
				})
			},
			expectedStatus:   http.StatusConflict,
			expectedUsername: "",
		},
		{
			name: "Bad data",
			req:  httptest.NewRequest("POST", "/register", bytes.NewReader([]byte{})),
			setup: func() {
			},
			expectedStatus:   http.StatusBadRequest,
			expectedUsername: "",
		},
		{
			name: "Wrong method",
			req:  httptest.NewRequest("GET", "/register", bytes.NewReader([]byte{})),
			setup: func() {
			},
			expectedStatus:   http.StatusBadRequest,
			expectedUsername: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			w := httptest.NewRecorder()

			handler.Register(w, tc.req)

			resp := w.Result()
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			var user users.User
			json.NewDecoder(resp.Body).Decode(&user)
			assert.Equal(t, tc.expectedUsername, user.Username)
		})
	}
}

func TestLogin(t *testing.T) {
	handler := setupTest()
	newUser := users.User{
		Username: "new_user",
		Password: "password123",
	}

	wrongUser := users.User{
		Username: "new_user",
		Password: "password124",
	}

	body, _ := json.Marshal(newUser)
	wrongData, _ := json.Marshal(wrongUser)

	testCases := []testCase{
		{
			name: "Valid Input",
			req:  httptest.NewRequest("POST", "/login", bytes.NewReader(body)),
			setup: func() {
				handler.UserDB.AddUser(&newUser)
			},
			expectedStatus:   http.StatusOK,
			expectedUsername: newUser.Username,
		},
		{
			name: "Wrong password",
			req:  httptest.NewRequest("POST", "/login", bytes.NewReader(wrongData)),
			setup: func() {
				handler.UserDB.AddUser(&newUser)
			},
			expectedStatus:   http.StatusUnauthorized,
			expectedUsername: "",
		},
		{
			name: "Bad data",
			req:  httptest.NewRequest("POST", "/login", bytes.NewReader([]byte{})),
			setup: func() {
			},
			expectedStatus:   http.StatusBadRequest,
			expectedUsername: "",
		},
		{
			name: "Wrong method",
			req:  httptest.NewRequest("GET", "/login", bytes.NewReader([]byte{})),
			setup: func() {
			},
			expectedStatus:   http.StatusBadRequest,
			expectedUsername: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.setup()
			w := httptest.NewRecorder()

			handler.Login(w, tc.req)

			resp := w.Result()
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			var user users.User
			json.NewDecoder(resp.Body).Decode(&user)
			assert.Equal(t, tc.expectedUsername, user.Username)
		})
	}
}

func TestCheckSession(t *testing.T) {
	handler := setupTest()
	newUser := users.User{
		Username: "new_user",
		Password: "password123",
	}
	handler.UserDB.AddUser(&newUser)
	session := handler.SessionDb.CreateSession(newUser.Username)

	testCases := []testCase{
		{
			name: "Valid input",
			req: func() *http.Request {
				req := httptest.NewRequest("GET", "/session", nil)
				req.AddCookie(&http.Cookie{
					Name:  SessionToken,
					Value: session.Token,
				})
				return req
			}(),
			expectedStatus:   http.StatusOK,
			expectedUsername: newUser.Username,
		},
		{
			name:             "Invalid cookie",
			req:              httptest.NewRequest("GET", "/session", nil),
			expectedStatus:   http.StatusUnauthorized,
			expectedUsername: "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()

			handler.CheckSession(w, tc.req)

			resp := w.Result()
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			var user users.User
			json.NewDecoder(resp.Body).Decode(&user)
			assert.Equal(t, tc.expectedUsername, user.Username)
		})
	}
}

func TestLogout(t *testing.T) {
	handler := setupTest()
	newUser := users.User{
		Username: "new_user",
		Password: "password123",
	}

	handler.UserDB.AddUser(&newUser)
	session := handler.SessionDb.CreateSession(newUser.Username)

	testCases := []testCase{
		{
			name: "Valid input",
			req: func() *http.Request {
				req := httptest.NewRequest("GET", "/session", nil)
				req.AddCookie(&http.Cookie{
					Name:  SessionToken,
					Value: session.Token,
				})
				return req
			}(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid cookie",
			req:            httptest.NewRequest("GET", "/session", nil),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()

			handler.Logout(w, tc.req)

			resp := w.Result()
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}
