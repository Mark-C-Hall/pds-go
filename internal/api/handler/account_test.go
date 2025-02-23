package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleCreateAccount(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		requestBody    interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "valid request",
			method: http.MethodPost,
			requestBody: CreateAccountRequest{
				Handle:   "test.bsky.social",
				Email:    "test@example.com",
				Password: "strongpassword123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "wrong HTTP method",
			method: http.MethodGet,
			requestBody: CreateAccountRequest{
				Handle: "test.bsky.social",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Method Not Supported",
		},
		{
			name:   "missing content-type header",
			method: http.MethodPost,
			requestBody: CreateAccountRequest{
				Handle: "test.bsky.social",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Content-Type must be application/json",
		},
		{
			name:   "missing handle",
			method: http.MethodPost,
			requestBody: CreateAccountRequest{
				Email:    "test@example.com",
				Password: "strongpassword123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Handle is required",
		},
		{
			name:   "missing password",
			method: http.MethodPost,
			requestBody: CreateAccountRequest{
				Handle: "test.bsky.social",
				Email:  "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Password is required",
		},
		{
			name:   "missing email",
			method: http.MethodPost,
			requestBody: CreateAccountRequest{
				Handle:   "test.bsky.social",
				Password: "strongpassword123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Email is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create body
			var body bytes.Buffer
			if str, ok := tt.requestBody.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.requestBody)
			}

			// Create test request
			req := httptest.NewRequest(tt.method, "/xrpc/com.atproto.server.createAccount", &body)

			// Only set Content-Type for test cases that should have it
			if tt.name != "missing content-type header" {
				req.Header.Set("Content-Type", "application/json")
			}

			// Create response recorder
			rec := httptest.NewRecorder()

			// Call handler
			HandleCreateAccount(rec, req)

			// Check status code
			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			// Verify error message if applicable
			if tt.expectedError != "" {
				var errorResp ErrorResponse
				if err := json.NewDecoder(rec.Body).Decode(&errorResp); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}
				if msg := errorResp.Message; msg != tt.expectedError {
					t.Errorf("expected error message %q, got %q", tt.expectedError, msg)
				}
			}
		})
	}
}
