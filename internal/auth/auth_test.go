package auth

import (
	"bytes"
	"encoding/json"
	"kudago/internal/users"
	"kudago/session"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCaseRegister struct {
	name             string
	req              *http.Request
	setupHandler     func()
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
		Username:    "new_user",
		Email:       "new_user@example.com",
		DateOfBirth: "1990-01-01",
		Password:    "password123",
	}

	body, _ := json.Marshal(newUser)

	testCases := []testCaseRegister{
		{
			name:             "Valid Input",
			req:              httptest.NewRequest("POST", "/register", bytes.NewReader(body)),
			setupHandler: func() {
				handler = setupTest() 
			},
			expectedStatus:   http.StatusCreated,
			expectedUsername: newUser.Username,
		},
		{
			name: "User already exists",
			req:  httptest.NewRequest("POST", "/register", bytes.NewReader(body)),
			setupHandler: func() {
				handler = setupTest() 
				handler.UserDB.AddUser(&newUser)
			},
			expectedStatus:   http.StatusConflict,
			expectedUsername: "",
		},
		{
			name: "Bad data",
			req:  httptest.NewRequest("POST", "/register", bytes.NewReader([]byte{})),
			setupHandler: func() {
				handler = setupTest() 
			},
			expectedStatus:   http.StatusBadRequest,
			expectedUsername: "",
		},
		{
			name: "Wrong method",
			req:  httptest.NewRequest("GET", "/register", bytes.NewReader([]byte{})),
			setupHandler: func() {
				handler = setupTest() 
			},
			expectedStatus:   http.StatusBadRequest,
			expectedUsername: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupHandler() 
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

	testCases := []testCaseRegister{
		{
			name: "Valid Input",
			req:  httptest.NewRequest("POST", "/login", bytes.NewReader(body)),
			setupHandler: func() {
				handler.UserDB.AddUser(&newUser)
			},
			expectedStatus:   http.StatusOK,
			expectedUsername: newUser.Username,
		},
		{
			name: "Wrong password",
			req:  httptest.NewRequest("POST", "/login", bytes.NewReader(wrongData)),
			setupHandler: func() {
				handler.UserDB.AddUser(&newUser)
			},
			expectedStatus:   http.StatusUnauthorized,
			expectedUsername: "",
		},
		{
			name: "Bad data",
			req:  httptest.NewRequest("POST", "/login", bytes.NewReader([]byte{})),
			setupHandler: func() {
			},
			expectedStatus:   http.StatusBadRequest,
			expectedUsername: "",
		},
		{
			name: "Wrong method",
			req:  httptest.NewRequest("GET", "/login", bytes.NewReader([]byte{})),
			setupHandler: func() {
			},
			expectedStatus:   http.StatusBadRequest,
			expectedUsername: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.setupHandler()
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
