package application_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hxzzf/calc_go/internal/application"
)

func TestHandleCalculate(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		requestBody    interface{}
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "valid expression",
			method: http.MethodPost,
			requestBody: map[string]string{
				"expression": "2 + 2",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"result": float64(4),
			},
		},
		{
			name:   "complex expression",
			method: http.MethodPost,
			requestBody: map[string]string{
				"expression": "(2 + 3) * 4",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"result": float64(20),
			},
		},
		{
			name:   "empty expression",
			method: http.MethodPost,
			requestBody: map[string]string{
				"expression": "",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody: map[string]interface{}{
				"error": "Expression cannot be empty",
			},
		},
		{
			name:   "invalid expression with letters",
			method: http.MethodPost,
			requestBody: map[string]string{
				"expression": "2 + a",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody: map[string]interface{}{
				"error": "Expression is not valid",
			},
		},
		{
			name:   "division by zero",
			method: http.MethodPost,
			requestBody: map[string]string{
				"expression": "1 / 0",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody: map[string]interface{}{
				"error": "Division by zero is not allowed",
			},
		},
		{
			name:   "consecutive operators",
			method: http.MethodPost,
			requestBody: map[string]string{
				"expression": "2 ++ 2",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody: map[string]interface{}{
				"error": "Consecutive operators are not allowed",
			},
		},
		{
			name:   "mismatched parentheses",
			method: http.MethodPost,
			requestBody: map[string]string{
				"expression": "(2 + 3",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody: map[string]interface{}{
				"error": "Parentheses are mismatched",
			},
		},
		{
			name:   "invalid request body",
			method: http.MethodPost,
			requestBody: map[string]string{
				"wrong_field": "2 + 2",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody: map[string]interface{}{
				"error": "Invalid request body",
			},
		},
		{
			name:           "wrong method",
			method:         http.MethodGet,
			requestBody:    nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody: map[string]interface{}{
				"error": "Method not allowed",
			},
		},
		{
			name:   "malformed JSON",
			method: http.MethodPost,
			requestBody: `{
				"expression": "2 + 2"
				invalid json
			}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody: map[string]interface{}{
				"error": "Invalid request body",
			},
		},
		{
			name:   "internal server error",
			method: http.MethodPost,
			requestBody: map[string]string{
				"expression": "1,7976931348623157 * 2",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Internal server error",
			},
		},
	}

	app := application.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error

			if tt.requestBody != nil {
				switch v := tt.requestBody.(type) {
				case string:
					body = []byte(v)
				default:
					body, err = json.Marshal(tt.requestBody)
					if err != nil {
						t.Fatalf("Failed to marshal request body: %v", err)
					}
				}
			}

			req := httptest.NewRequest(tt.method, "/api/v1/calculate", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			http.HandlerFunc(app.HandleCalculate).ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]interface{}
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			for key, expectedValue := range tt.expectedBody {
				if actualValue, ok := response[key]; !ok {
					t.Errorf("Expected response to contain key %q", key)
				} else if expectedValue != actualValue {
					t.Errorf("Expected %q to be %v, got %v", key, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestConfigFromEnv(t *testing.T) {
	originalPort := os.Getenv("PORT")
	defer os.Setenv("PORT", originalPort)

	tests := []struct {
		name     string
		portEnv  string
		expected string
	}{
		{
			name:     "default port",
			portEnv:  "",
			expected: "8080",
		},
		{
			name:     "custom port",
			portEnv:  "3000",
			expected: "3000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("PORT", tt.portEnv)
			config := application.ConfigFromEnv()
			if config.Port != tt.expected {
				t.Errorf("Expected port %s, got %s", tt.expected, config.Port)
			}
		})
	}
}
