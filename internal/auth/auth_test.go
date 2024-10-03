package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"kudago/internal/users"
	"kudago/session"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name              string
	req               *http.Request
	setup             func()
	expectedUsername  string
	expectedStatus    int
	expectedErrStatus int
}

func setupTest() *Handler {
	handler := &Handler{
		UserDB:    *users.NewUserDB(),
		SessionDb: *session.NewSessionDB(),
	}
	return handler
}

func TestRegister(t *testing.T) {
	t.Parallel()
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
			expectedStatus:    http.StatusCreated,
			expectedErrStatus: 0,
			expectedUsername:  newUser.Username,
		},
		{
			name: "User already exists",
			req:  httptest.NewRequest("POST", "/register", bytes.NewReader(body)),
			setup: func() {
				handler.UserDB.AddUser(&newUser)
			},
			expectedStatus:    http.StatusOK,
			expectedErrStatus: http.StatusConflict,
			expectedUsername:  "",
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
			expectedStatus:    http.StatusOK,
			expectedErrStatus: http.StatusConflict,
			expectedUsername:  "",
		},
		{
			name: "Bad data",
			req:  httptest.NewRequest("POST", "/register", bytes.NewReader([]byte{})),
			setup: func() {
			},
			expectedStatus:    http.StatusOK,
			expectedErrStatus: http.StatusBadRequest,
			expectedUsername:  "",
		},
		{
			name: "Wrong method",
			req:  httptest.NewRequest("GET", "/register", bytes.NewReader([]byte{})),
			setup: func() {
			},
			expectedStatus:    http.StatusOK,
			expectedErrStatus: http.StatusBadRequest,
			expectedUsername:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			w := httptest.NewRecorder()

			handler.Register(w, tc.req)

			resp := w.Result()
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			body, _ := io.ReadAll(resp.Body)

			var responseErr authError
			json.Unmarshal(body, &responseErr)
			assert.Equal(t, tc.expectedErrStatus, responseErr.Code)

			var user users.User
			json.Unmarshal(body, &user)
			assert.Equal(t, tc.expectedUsername, user.Username)
		})
	}
}

func TestLogin(t *testing.T) {
	t.Parallel()
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
			expectedStatus:    http.StatusOK,
			expectedErrStatus: 0,
			expectedUsername:  newUser.Username,
		},
		{
			name: "Wrong password",
			req:  httptest.NewRequest("POST", "/login", bytes.NewReader(wrongData)),
			setup: func() {
				handler.UserDB.AddUser(&newUser)
			},
			expectedStatus:    http.StatusOK,
			expectedErrStatus: http.StatusUnauthorized,
			expectedUsername:  "",
		},
		{
			name: "Bad data",
			req:  httptest.NewRequest("POST", "/login", bytes.NewReader([]byte{})),
			setup: func() {
			},
			expectedStatus:    http.StatusOK,
			expectedErrStatus: http.StatusBadRequest,
			expectedUsername:  "",
		},
		{
			name: "Wrong method",
			req:  httptest.NewRequest("GET", "/login", bytes.NewReader([]byte{})),
			setup: func() {
			},
			expectedStatus:    http.StatusOK,
			expectedErrStatus: http.StatusBadRequest,
			expectedUsername:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			tc.setup()
			w := httptest.NewRecorder()

			handler.Login(w, tc.req)

			resp := w.Result()
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			body, _ := io.ReadAll(resp.Body)

			var responseErr authError
			json.Unmarshal(body, &responseErr)
			assert.Equal(t, tc.expectedErrStatus, responseErr.Code)

			var user users.User
			json.Unmarshal(body, &user)
			assert.Equal(t, tc.expectedUsername, user.Username)
		})
	}
}

func TestCheckSession(t *testing.T) {
	t.Parallel()
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
			setup:             func() {},
			expectedStatus:    http.StatusOK,
			expectedErrStatus: 0,
			expectedUsername:  newUser.Username,
		},
		{
			name:              "Invalid cookie",
			req:               httptest.NewRequest("GET", "/session", nil),
			expectedErrStatus: http.StatusUnauthorized,
			expectedStatus:    http.StatusOK,
			setup:             func() {},
			expectedUsername:  "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			w := httptest.NewRecorder()

			handler.CheckSession(w, tc.req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			var user users.User
			json.Unmarshal(body, &user)
			assert.Equal(t, tc.expectedUsername, user.Username)

			var responseErr authError
			json.Unmarshal(body, &responseErr)
			assert.Equal(t, tc.expectedErrStatus, responseErr.Code)
		})
	}

}

func TestLogout(t *testing.T) {
	t.Parallel()

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
			expectedErrStatus: 0,
			expectedStatus:    http.StatusOK,
		},
		{
			name:              "Invalid cookie",
			req:               httptest.NewRequest("GET", "/session", nil),
			expectedErrStatus: http.StatusUnauthorized,
			expectedStatus:    http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()

			handler.Logout(w, tc.req)

			resp := w.Result()
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			body, _ := io.ReadAll(resp.Body)
			var responseErr authError
			json.Unmarshal(body, &responseErr)

			assert.Equal(t, tc.expectedErrStatus, responseErr.Code)

		})
	}
}
