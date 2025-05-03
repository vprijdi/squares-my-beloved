package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestUserRegistration(t *testing.T) {
	db, cleanup := newInMemTestDatabase(t)
	defer cleanup()

	testApp := newTestApplication(t, db)
	testMux := testApp.Mount()

	t.Run("should return 405 for wrong method", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/auth/register", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, testMux)

		checkResponseCode(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("should return 400 for empty body", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBuffer([]byte{}))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := executeRequest(req, testMux)
		checkResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return 400 for invalid payload", func(t *testing.T) {
		invalidPayloads := []struct {
			name    string
			payload map[string]interface{}
		}{
			{
				name: "missing email",
				payload: map[string]interface{}{
					"username": "testuser",
					"password": "validPass123!",
				},
			},
			{
				name: "invalid email",
				payload: map[string]interface{}{
					"email":    "notanemail",
					"username": "testuser",
					"password": "validPass123!",
				},
			},
			{
				name: "short password",
				payload: map[string]interface{}{
					"email":    "test@example.com",
					"username": "testuser",
					"password": "short",
				},
			},
			{
				name: "missing username",
				payload: map[string]interface{}{
					"email":    "test@example.com",
					"password": "validPass123!",
				},
			},
		}

		for _, tc := range invalidPayloads {
			t.Run(tc.name, func(t *testing.T) {
				body, _ := json.Marshal(tc.payload)
				req, err := http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBuffer(body))
				if err != nil {
					t.Fatal(err)
				}

				rr := executeRequest(req, testMux)
				checkResponseCode(t, http.StatusBadRequest, rr.Code)
			})
		}
	})

	t.Run("should return 400 for duplicate email or username with proper error structure", func(t *testing.T) {
		// First create a user
		initialUser := map[string]interface{}{
			"email":    "existing@example.com",
			"username": "existinguser",
			"password": "validPass123!",
		}

		body, _ := json.Marshal(initialUser)
		req, _ := http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		executeRequest(req, testMux)

		testCases := []struct {
			name        string
			payload     map[string]interface{}
			expectedKey string
			expectedMsg string
		}{
			{
				name: "duplicate email",
				payload: map[string]interface{}{
					"email":    "existing@example.com",
					"username": "newuser",
					"password": "validPass123!",
				},
				expectedKey: "email",
				expectedMsg: "a user with that email already exist",
			},
			{
				name: "duplicate username",
				payload: map[string]interface{}{
					"email":    "new@example.com",
					"username": "existinguser",
					"password": "validPass123!",
				},
				expectedKey: "username",
				expectedMsg: "a user with that username already exist",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				body, _ := json.Marshal(tc.payload)
				req, _ := http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")

				rr := executeRequest(req, testMux)
				checkResponseCode(t, http.StatusBadRequest, rr.Code)

				// Verify error response structure
				var response map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Fatalf("could not decode response: %v", err)
				}

				// Check for expected error fields
				if _, ok := response["error"]; !ok {
					t.Error("expected 'error' field in response")
				}

				// For more detailed validation if your error has a specific structure
				if errorVal, ok := response["error"].(map[string]interface{}); ok {
					if msg, exists := errorVal[tc.expectedKey]; !exists {
						t.Errorf("expected error to contain key %q", tc.expectedKey)
					} else if msg != tc.expectedMsg {
						t.Errorf("expected error message %q, got %q", tc.expectedMsg, msg)
					}
				} else if errorMsg, ok := response["error"].(string); ok {
					if !strings.Contains(errorMsg, tc.expectedMsg) {
						t.Errorf("expected error message to contain %q, got %q", tc.expectedMsg, errorMsg)
					}
				} else {
					t.Error("unexpected error response format")
				}
			})
		}
	})

	t.Run("should successfully register a new user", func(t *testing.T) {
		payload := map[string]interface{}{
			"email":    "test@example.com",
			"username": "testuser",
			"password": "validPass123!",
		}

		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := executeRequest(req, testMux)
		checkResponseCode(t, http.StatusCreated, rr.Code)

		// Verify response structure matches actual API response
		var response struct {
			Data struct {
				ID          int         `json:"id"`
				Email       string      `json:"email"`
				Username    string      `json:"username"`
				DisplayName interface{} `json:"display_name"`
				IsActive    bool        `json:"is_active"`
				CreatedAt   string      `json:"created_at"`
				Token       string      `json:"token"`
			} `json:"data"`
		}

		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		// 1. Verify token format (UUID)
		_, err = uuid.Parse(response.Data.Token)
		if err != nil {
			t.Errorf("token is not a valid UUID: %v", err)
		}

		// Now check the values at the correct nesting level
		if response.Data.Username != payload["username"] {
			t.Errorf("expected username %s, got %s", payload["username"], response.Data.Username)
		}

		if response.Data.Email != payload["email"] {
			t.Errorf("expected email %s, got %s", payload["email"], response.Data.Email)
		}

		if response.Data.Token == "" {
			t.Error("expected a token in the response")
		}
	})
}
