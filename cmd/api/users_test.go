package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/muzhiknastya/squares-my-beloved/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestGetUserHandler(t *testing.T) {
	db, cleanup := newInMemTestDatabase(t)
	defer cleanup()

	testApp := newTestApplication(t, db)
	testMux := testApp.Mount()

	activeUser := &store.User{
		Email:    "activeuser@example.com",
		Username: "activeuser",
		IsActive: true,
	}
	activeUser.Password.Set(testPassword)
	activeUserID := createTestUser(t, db, testApp, activeUser)
	validToken := getAuthToken(t, testApp, testMux, activeUser.Email, testPassword)

	inactiveUser := &store.User{
		Email:    "inactiveuser@example.com",
		Username: "inactiveuser",
		IsActive: false,
	}
	inactiveUser.Password.Set(testPassword)
	_ = createTestUser(t, db, testApp, inactiveUser)

	tests := []struct {
		name           string
		userID         string
		token          string
		wantStatus     int
		wantErrorKey   string
		wantErrorValue string
	}{
		{
			name:       "successfully get active user",
			userID:     fmt.Sprintf("%d", activeUserID),
			token:      validToken,
			wantStatus: http.StatusOK,
		},
		{
			name:       "inactive user still returns 200",
			userID:     fmt.Sprintf("%d", activeUserID),
			token:      validToken,
			wantStatus: http.StatusOK,
		},
		{
			name:           "unauthorized access returns 401",
			userID:         fmt.Sprintf("%d", activeUserID),
			token:          "invalidtoken",
			wantStatus:     http.StatusUnauthorized,
			wantErrorKey:   "error",
			wantErrorValue: "unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/users/%s", tt.userID), nil)
			req.Header.Set("Authorization", "Bearer "+tt.token)

			rr := executeRequest(req, testMux)
			checkResponseCode(t, tt.wantStatus, rr.Code)

			if tt.wantStatus == http.StatusOK {
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

				// Verify response fields
				assert.Equal(t, activeUserID, response.Data.ID, "Expected user ID %d, got %d", activeUserID, response.Data.ID)
				assert.Equal(t, activeUser.Email, response.Data.Email, "Expected email %s, got %s", activeUser.Email, response.Data.Email)
				assert.Equal(t, activeUser.Username, response.Data.Username, "Expected username %s, got %s", activeUser.Username, response.Data.Username)
				assert.Nil(t, response.Data.DisplayName, "Expected display_name to be nil")
				if response.Data.CreatedAt == "" {
					t.Error("Expected created_at to be set")
				}

				// Verify created_at timestamp format
				if _, err := time.Parse(time.RFC3339, response.Data.CreatedAt); err != nil {
					t.Errorf("Invalid created_at format: %v", err)
				}
			} else {
				var errorResponse map[string]interface{}
				if err := json.NewDecoder(rr.Body).Decode(&errorResponse); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}

				errorMsg, ok := errorResponse[tt.wantErrorKey].(string)
				if !ok {
					t.Errorf("Expected error key '%s' in response", tt.wantErrorKey)
				} else if !strings.Contains(errorMsg, tt.wantErrorValue) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.wantErrorValue, errorMsg)
				}
			}
		})
	}
}
