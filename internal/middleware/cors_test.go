package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSMiddleware(t *testing.T) {
	// Вспомогательная функция для тестирования middleware
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что middleware передал запрос в следующий обработчик
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name               string
		originHeader       string
		expectedStatusCode int
		expectedHeaders    map[string]string
	}{
		{
			name:               "Valid origin - localhost",
			originHeader:       "http://localhost",
			expectedStatusCode: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "http://localhost",
				"Access-Control-Allow-Methods":     "GET, POST, OPTIONS, PUT, DELETE",
				"Access-Control-Allow-Headers":     "Content-Type",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			name:               "Valid origin - vyhodnoy.online",
			originHeader:       "http://vyhodnoy.online",
			expectedStatusCode: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "http://vyhodnoy.online",
				"Access-Control-Allow-Methods":     "GET, POST, OPTIONS, PUT, DELETE",
				"Access-Control-Allow-Headers":     "Content-Type",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			name:               "Invalid origin",
			originHeader:       "http://unauthorized-origin.com",
			expectedStatusCode: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Methods":     "GET, POST, OPTIONS, PUT, DELETE",
				"Access-Control-Allow-Headers":     "Content-Type",
				"Access-Control-Allow-Credentials": "true",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем запрос с необходимыми заголовками
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Origin", tt.originHeader)

			// Создаем новый тестовый сервер с middleware
			rr := httptest.NewRecorder()
			handler := CORSMiddleware(testHandler)

			// Выполняем запрос
			handler.ServeHTTP(rr, req)

			// Проверяем статус код
			if rr.Code != tt.expectedStatusCode {
				t.Errorf("expected status code %v, got %v", tt.expectedStatusCode, rr.Code)
			}

			// Проверяем заголовки
			for key, expectedValue := range tt.expectedHeaders {
				actualValue := rr.Header().Get(key)
				if actualValue != expectedValue {
					t.Errorf("expected header %v: %v, got %v", key, expectedValue, actualValue)
				}
			}
		})
	}
}
