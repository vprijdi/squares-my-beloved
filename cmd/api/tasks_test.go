package api

import (
	"context"
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

	// Setup test environment
	userID := createTestUser(t, db, testApp)
	validToken := getAuthToken(t, testApp, testMux)

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
			_, err := db.Exec("DELETE FROM tasks")
			if err != nil {
				t.Fatalf("Failed to clear tasks: %v", err)
			}

			req, _ := http.NewRequest("POST", "/v1/tasks", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+validToken)

			rr := httptest.NewRecorder()
			testMux.ServeHTTP(rr, req)

			checkResponseCode(t, tt.wantStatus, rr.Code)

			if tt.wantStatus == http.StatusCreated {
				var dbCount int
				err := db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM tasks WHERE user_id = $1", userID).Scan(&dbCount)
				if err != nil {
					t.Fatalf("Failed to count tasks: %v", err)
				}
				if dbCount != tt.wantTaskCount {
					t.Errorf("Expected %d tasks in DB, got %d", tt.wantTaskCount, dbCount)
				}
			}

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

	t.Run("unauthorized access", func(t *testing.T) {
		payload := `[{"title": "Task", "is_optional": false, "tier": 1}]`
		req, _ := http.NewRequest("POST", "/v1/tasks", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer invalidtoken")

		rr := httptest.NewRecorder()
		testMux.ServeHTTP(rr, req)

		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})
}

func TestGetTaskHandler(t *testing.T) {
	db, cleanup := newInMemTestDatabase(t)
	defer cleanup()

	testApp := newTestApplication(t, db)
	testMux := testApp.Mount()

	userID := createTestUser(t, db, testApp)
	task := createTestTask(t, context.Background(), testApp, userID)
	validToken := getAuthToken(t, testApp, testMux)

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

		// Verify response fields
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

		// Verify timestamp formats
		if _, err := time.Parse(time.RFC3339, response.Data.CreatedAt); err != nil {
			t.Errorf("Invalid created_at format: %v", err)
		}
		if _, err := time.Parse(time.RFC3339, response.Data.UpdatedAt); err != nil {
			t.Errorf("Invalid updated_at format: %v", err)
		}
	})

	t.Run("nonexistent task returns 404", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/tasks/999999", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)

		rr := httptest.NewRecorder()
		testMux.ServeHTTP(rr, req)

		checkResponseCode(t, http.StatusNotFound, rr.Code)
	})

	t.Run("invalid task ID format returns 500", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/tasks/invalid", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)

		rr := httptest.NewRecorder()
		testMux.ServeHTTP(rr, req)

		checkResponseCode(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("unauthorized access returns 401", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/tasks/%d", task.ID), nil)
		req.Header.Set("Authorization", "Bearer invalidtoken")

		rr := httptest.NewRecorder()
		testMux.ServeHTTP(rr, req)

		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})
}

func TestCompleteTaskHandler(t *testing.T) {
	db, cleanup := newInMemTestDatabase(t)
	defer cleanup()

	testApp := newTestApplication(t, db)
	testMux := testApp.Mount()

	// Create test user
	userID := createTestUser(t, db, testApp)

	// Create test task
	task := &store.Task{
		Title:      "Test Task",
		IsOptional: false,
		Tier:       1,
		UserID:     userID,
	}
	if err := testApp.store.Tasks.CreateTasks(context.Background(), []*store.Task{task}); err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	// Get valid token
	validToken := getAuthToken(t, testApp, testMux)

	t.Run("successfully complete task", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", fmt.Sprintf("/v1/tasks/%d/complete", task.ID), nil)
		req.Header.Set("Authorization", "Bearer "+validToken)

		rr := httptest.NewRecorder()
		testMux.ServeHTTP(rr, req)

		checkResponseCode(t, http.StatusNoContent, rr.Code)

		var completed bool
		err := db.QueryRowContext(context.Background(),
			"SELECT completed FROM tasks WHERE id = $1", task.ID).Scan(&completed)
		if err != nil {
			t.Fatalf("Failed to query task: %v", err)
		}
		if !completed {
			t.Error("Expected task to be marked as completed")
		}
	})

	t.Run("nonexistent task returns 404", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", "/v1/tasks/999999/complete", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)

		rr := httptest.NewRecorder()
		testMux.ServeHTTP(rr, req)

		checkResponseCode(t, http.StatusNotFound, rr.Code)
	})

	t.Run("unauthorized access returns 401", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", fmt.Sprintf("/v1/tasks/%d/complete", task.ID), nil)
		req.Header.Set("Authorization", "Bearer invalidtoken")

		rr := httptest.NewRecorder()
		testMux.ServeHTTP(rr, req)

		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("invalid task ID format returns 500", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", "/v1/tasks/not-a-number/complete", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)

		rr := httptest.NewRecorder()
		testMux.ServeHTTP(rr, req)

		checkResponseCode(t, http.StatusInternalServerError, rr.Code)
	})
}
