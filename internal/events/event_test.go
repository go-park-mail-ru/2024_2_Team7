package events

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupTest() *Handler {
	return &Handler{
		EventDB: *NewEventDB(),
	}
}

func TestGetAllEvents(t *testing.T) {
t.Parallel()

	handler := setupTest()
	req := httptest.NewRequest("GET", "/events", nil)
	w := httptest.NewRecorder()

	handler.GetAllEvents(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var events []Event
	err := json.NewDecoder(resp.Body).Decode(&events)
	assert.NoError(t, err)

	expectedEvents := handler.EventDB.Events
	assert.Equal(t, len(expectedEvents), len(events))
}

func TestGetEventsByTag(t *testing.T) {
	handler := setupTest()

	testEvent := Event{
		ID:    89,
		Title: "Test",
		Tag:   []string{"test"},
	}

	handler.EventDB.Events = append(handler.EventDB.Events, testEvent)

	testCases := []struct {
		name           string
		tag            string
		expectedStatus int
		expectedCount  []Event
	}{
		{
			name:           "Get by festival tag events",
			tag:            "test",
			expectedStatus: http.StatusOK,
			expectedCount:  []Event{testEvent},
		},
		{
			name:           "Tag not found",
			tag:            "unknown",
			expectedStatus: http.StatusNoContent,
			expectedCount:  []Event{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/events/"+tc.tag, nil)
			req = mux.SetURLVars(req, map[string]string{"tag": tc.tag})
			w := httptest.NewRecorder()

			handler.GetEventsByTag(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedStatus == http.StatusOK {
				var events []Event
				err := json.NewDecoder(resp.Body).Decode(&events)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCount, events)
			}
		})
	}
}
