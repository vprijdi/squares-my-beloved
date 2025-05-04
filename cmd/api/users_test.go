package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetUserHandler(t *testing.T) {
	db, cleanup := newInMemTestDatabase(t)
	defer cleanup()

	testApp := newTestApplication(t, db)
	testMux := testApp.Mount()

	t.Run("successfully get active user", func(t *testing.T) {
		cleanUserTable(t, db)
		userID := createTestUser(t, db, testApp, true)

		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/users/%d", userID), nil)
		rr := httptest.NewRecorder()
		testMux.ServeHTTP(rr, req)

		checkResponseCode(t, http.StatusOK, rr.Code)

		var response struct {
			Data struct {
				ID          int64   `json:"id"`
				Email       string  `json:"email"`
				Username    string  `json:"username"`
				DisplayName *string `json:"display_name"`
				IsActive    bool    `json:"is_active"`
				CreatedAt   string  `json:"created_at"`
			} `json:"data"`
		}

		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Verify response structure
		if response.Data.ID != userID {
			t.Errorf("Expected user ID %d, got %d", userID, response.Data.ID)
		}
		if response.Data.Email == "" {
			t.Error("Expected email to be set")
		}
		if response.Data.Username == "" {
			t.Error("Expected username to be set")
		}
		if response.Data.DisplayName != nil {
			t.Error("Expected display_name to be null")
		}
		if !response.Data.IsActive {
			t.Error("Expected user to be active")
		}
		if response.Data.CreatedAt == "" {
			t.Error("Expected created_at to be set")
		}

		// Verify timestamp format
		if _, err := time.Parse(time.RFC3339, response.Data.CreatedAt); err != nil {
			t.Errorf("Invalid created_at format: %v", err)
		}
	})

	t.Run("inactive user returns 404", func(t *testing.T) {
		cleanUserTable(t, db)
		userID := createTestUser(t, db, testApp, false)

		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/users/%d", userID), nil)
		rr := httptest.NewRecorder()
		testMux.ServeHTTP(rr, req)

		checkResponseCode(t, http.StatusNotFound, rr.Code)

		var response map[string]string
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}
		if response["error"] == "" {
			t.Error("Expected error message in response")
		}
	})

	t.Run("nonexistent user returns 404", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/users/999999", nil)
		rr := httptest.NewRecorder()
		testMux.ServeHTTP(rr, req)

		checkResponseCode(t, http.StatusNotFound, rr.Code)
	})

	t.Run("invalid user ID format returns 500", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/users/invalid", nil)
		rr := httptest.NewRecorder()
		testMux.ServeHTTP(rr, req)

		checkResponseCode(t, http.StatusInternalServerError, rr.Code)
	})
}
