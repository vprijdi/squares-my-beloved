package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/muzhiknastya/squares-my-beloved/internal/store"
)

func TestCreateBatchTasksHandler(t *testing.T) {
	db, cleanup := newInMemTestDatabase(t)
	defer cleanup()

	testApp := newTestApplication(t, db)
	testMux := testApp.Mount()

	// Create test user
	ctx := context.Background()
	var userID int64
	err := store.WithTx(db, ctx, func(tx *sql.Tx) error {
		user := &store.User{
			Email:    "testuser@example.com",
			Username: "testuser",
			IsActive: true,
		}
		if err := user.Password.Set("testpassword"); err != nil {
			return err
		}
		if err := testApp.store.Users.Create(ctx, tx, user); err != nil {
			return err
		}
		userID = user.ID
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Get valid token through login endpoint
	loginBody := fmt.Sprintf(`{"email":"testuser@example.com","password":"testpassword"}`)
	req, _ := http.NewRequest("POST", "/v1/auth/token", strings.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testMux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("Login failed: %d - %s", rr.Code, rr.Body.String())
	}

	var tokenResp struct {
		Data string `json:"data"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&tokenResp); err != nil {
		t.Fatalf("Failed to decode token: %v", err)
	}
	validToken := tokenResp.Data

	// Test cases
	tests := []struct {
		name           string
		payload        string
		wantStatus     int
		wantTaskCount  int
		wantErrorKey   string
		wantErrorValue string
	}{
		{
			name:          "successful batch creation",
			payload:       `[{"title": "Valid Task 1", "is_optional": false, "tier": 1}, {"title": "Valid Task 2", "is_optional": true, "tier": 2}]`,
			wantStatus:    http.StatusCreated,
			wantTaskCount: 2,
		},
		{
			name:           "empty task list",
			payload:        `[]`,
			wantStatus:     http.StatusBadRequest,
			wantErrorKey:   "error",
			wantErrorValue: "at least one task is required",
		},
		{
			name:           "empty title",
			payload:        `[{"title": "", "is_optional": false, "tier": 1}]`,
			wantStatus:     http.StatusBadRequest,
			wantErrorKey:   "error",
			wantErrorValue: "Title",
		},
		{
			name:           "invalid tier",
			payload:        `[{"title": "Task", "is_optional": false, "tier": 3}]`,
			wantStatus:     http.StatusBadRequest,
			wantErrorKey:   "error",
			wantErrorValue: "Tier",
		},
		{
			name:           "missing title",
			payload:        `[{"is_optional": false, "tier": 1}]`,
			wantStatus:     http.StatusBadRequest,
			wantErrorKey:   "error",
			wantErrorValue: "Title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear tasks before each test case
			_, err := db.Exec("DELETE FROM tasks")
			if err != nil {
				t.Fatalf("Failed to clear tasks: %v", err)
			}

			req, _ := http.NewRequest("POST", "/v1/tasks", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+validToken)

			rr := httptest.NewRecorder()
			testMux.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d: %s", tt.wantStatus, rr.Code, rr.Body.String())
			}

			// Verify successful response
			if tt.wantStatus == http.StatusCreated {

				// Verify tasks were created in database
				var dbCount int
				err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM tasks WHERE user_id = $1", userID).Scan(&dbCount)
				if err != nil {
					t.Fatalf("Failed to count tasks: %v", err)
				}
				if dbCount != tt.wantTaskCount {
					t.Errorf("Expected %d tasks in DB, got %d", tt.wantTaskCount, dbCount)
				}
			}

			// Verify error response
			if tt.wantStatus == http.StatusBadRequest {
				var response map[string]interface{}
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}

				errorMsg, ok := response[tt.wantErrorKey].(string)
				if !ok {
					t.Errorf("Expected error key '%s' in response", tt.wantErrorKey)
				} else if !strings.Contains(errorMsg, tt.wantErrorValue) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.wantErrorValue, errorMsg)
				}
			}
		})
	}

	// Test unauthorized access
	t.Run("unauthorized access", func(t *testing.T) {
		payload := `[{"title": "Task", "is_optional": false, "tier": 1}]`
		req, _ := http.NewRequest("POST", "/v1/tasks", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer invalidtoken")

		rr := httptest.NewRecorder()
		testMux.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rr.Code)
		}
	})
}
