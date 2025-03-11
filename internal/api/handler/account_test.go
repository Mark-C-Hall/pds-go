package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/mark-c-hall/pds-go/internal/api/httputil"
	"github.com/mark-c-hall/pds-go/internal/core/model"
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
			name: "valid request",
			requestBody: CreateAccountRequest{
				Handle:   "test.bsky.social",
				Email:    "test@example.com",
				Password: "strongpassword123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "missing handle",
			requestBody: CreateAccountRequest{
				Email:    "test@example.com",
				Password: "strongpassword123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Handle is required",
		},
		{
			name: "missing password",
			requestBody: CreateAccountRequest{
				Handle: "test.bsky.social",
				Email:  "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Password is required",
		},
		{
			name: "missing email",
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
			rec, req := setupTestRequest(http.MethodPost, "/xrpc/com.atproto.server.createAccount", tt.requestBody)

			// Create a mock service
			mockService := &MockAccountService{}

			// Create a test logger that discards output
			testLogger := log.New(io.Discard, "", 0)

			// Create the handler with mock dependencies
			handler := NewAccountHandler(mockService, testLogger)

			// Call the handler method
			handler.HandleCreateAccount(rec, req)

			// Check status code
			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			// Verify error message if applicable
			if tt.expectedError != "" {
				var errorResp httputil.ErrorResponse
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

func setupTestRequest(method, path string, body interface{}) (*httptest.ResponseRecorder, *http.Request) {
	// Create body
	var bodyBuffer bytes.Buffer
	if str, ok := body.(string); ok {
		bodyBuffer.WriteString(str)
	} else {
		json.NewEncoder(&bodyBuffer).Encode(body)
	}

	// Create test request
	req := httptest.NewRequest(method, path, &bodyBuffer)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rec := httptest.NewRecorder()

	return rec, req
}

// MockAccountService is a mock implementation of service.AccountService
type MockAccountService struct{}

// CreateAccount implements the AccountService interface for testing
func (m *MockAccountService) CreateAccount(ctx context.Context, handle syntax.Handle, email, password string) (*model.Account, error) {
	// For testing, return a mock account or error based on the input
	if handle == "" {
		return nil, fmt.Errorf("handle cannot be empty")
	}

	// Mock successful account creation
	return &model.Account{
		DID:       syntax.DID(fmt.Sprintf("did:plc:%s", handle)),
		Handle:    handle,
		Email:     email,
		CreatedAt: time.Now().UTC(),
	}, nil
}
