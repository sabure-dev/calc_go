package application

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		body          interface{}
		expectedCode  int
		expectedError bool
	}{
		{
			name:          "valid calculation",
			method:        http.MethodPost,
			body:          CalculationRequest{Expression: "2+2"},
			expectedCode:  http.StatusOK,
			expectedError: false,
		},
		{
			name:          "invalid method",
			method:        http.MethodGet,
			body:          nil,
			expectedCode:  http.StatusMethodNotAllowed,
			expectedError: true,
		},
		{
			name:          "invalid expression",
			method:        http.MethodPost,
			body:          CalculationRequest{Expression: "2++2"},
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
		{
			name:          "invalid json",
			method:        http.MethodPost,
			body:          "invalid json",
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
	}

	app := New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			var err error

			if tt.body != nil {
				reqBody, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatal(err)
				}
			}

			req := httptest.NewRequest(tt.method, "/calculate", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.calculateHandler)

			wrapped := &responseWriter{
				ResponseWriter: rr,
				status:         http.StatusOK,
			}

			handler.ServeHTTP(wrapped, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, rr.Code)
			}

			if tt.expectedError && wrapped.error == "" {
				t.Error("expected error message, got empty string")
			}

			if !tt.expectedError && wrapped.error != "" {
				t.Errorf("expected no error, got: %s", wrapped.error)
			}
		})
	}
}
