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
	"time"

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

func TestGetTaskHandler(t *testing.T) {
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

	// Create test task
	task := &store.Task{
		Title:      "Test Task",
		IsOptional: false,
		Tier:       1,
		UserID:     userID,
	}
	if err := testApp.store.Tasks.CreateTasks(ctx, []*store.Task{task}); err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	// Get valid token through login endpoint
	loginBody := fmt.Sprintf(`{"email":"testuser@example.com","password":"testpassword"}`)
	req, _ := http.NewRequest("POST", "/v1/auth/token", strings.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testMux.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusCreated, rr.Code)

	var tokenResp struct {
		Data string `json:"data"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&tokenResp); err != nil {
		t.Fatalf("Failed to decode token: %v", err)
	}
	validToken := tokenResp.Data

	t.Run("successfully get task", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/tasks/%d", task.ID), nil)
		req.Header.Set("Authorization", "Bearer "+validToken)

		rr := httptest.NewRecorder()
		testMux.ServeHTTP(rr, req)

		checkResponseCode(t, http.StatusOK, rr.Code)

		var response struct {
			Data struct {
				ID               int64  `json:"id"`
				Title            string `json:"title"`
				IsOptional       bool   `json:"is_optional"`
				Tier             int32  `json:"tier"`
				UserID           int64  `json:"user_id"`
				CreatedAt        string `json:"created_at"`
				UpdatedAt        string `json:"updated_at"`
				CompletionsCount int    `json:"completions_count"`
				IsCompleted      bool   `json:"is_completed"`
			} `json:"data"`
		}

		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Data.ID != task.ID {
			t.Errorf("Expected task ID %d, got %d", task.ID, response.Data.ID)
		}
		if response.Data.Title != task.Title {
			t.Errorf("Expected title '%s', got '%s'", task.Title, response.Data.Title)
		}
		if response.Data.IsOptional != task.IsOptional {
			t.Errorf("Expected is_optional %v, got %v", task.IsOptional, response.Data.IsOptional)
		}
		if response.Data.Tier != task.Tier {
			t.Errorf("Expected tier %d, got %d", task.Tier, response.Data.Tier)
		}
		if response.Data.UserID != userID {
			t.Errorf("Expected user ID %d, got %d", userID, response.Data.UserID)
		}
		if response.Data.CreatedAt == "" {
			t.Error("Expected created_at to be set")
		}
		if response.Data.UpdatedAt == "" {
			t.Error("Expected updated_at to be set")
		}
		if response.Data.CompletionsCount != 0 {
			t.Errorf("Expected completions_count 0, got %d", response.Data.CompletionsCount)
		}
		if response.Data.IsCompleted {
			t.Error("Expected is_completed to be false for new task")
		}

		// Verify createdAt is valid RFC3339 format
		if _, err := time.Parse(time.RFC3339, response.Data.CreatedAt); err != nil {
			t.Errorf("Invalid created_at format: %v", err)
		}
		// Verify updatedAt is valid RFC3339 format
		if _, err := time.Parse(time.RFC3339, response.Data.UpdatedAt); err != nil {
			t.Errorf("Invalid updated_at format: %v", err)
		}
	})

	t.Run("nonexistent task returns 404", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/tasks/999999", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)

		rr := httptest.NewRecorder()
		testMux.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", rr.Code)
		}

		var response map[string]string
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}
		if response["error"] == "" {
			t.Error("Expected error message in response")
		}
	})

	t.Run("invalid task ID format returns 500", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/tasks/invalid", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)

		rr := httptest.NewRecorder()
		testMux.ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})

	t.Run("unauthorized access returns 401", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/tasks/%d", task.ID), nil)
		req.Header.Set("Authorization", "Bearer invalidtoken")

		rr := httptest.NewRecorder()
		testMux.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rr.Code)
		}
	})
}
