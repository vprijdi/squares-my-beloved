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

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/muzhiknastya/squares-my-beloved/internal/auth"
	"github.com/muzhiknastya/squares-my-beloved/internal/services"
	"github.com/muzhiknastya/squares-my-beloved/internal/store"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, postgres *sql.DB) *Application {
	t.Helper()

	logger := zap.NewNop().Sugar()

	cfg := Config{
		Addr: "localhost:8080", // Not actually used for tests
		DB: DBConfig{
			Addr:         "localhost:5433",
			MaxOpenConns: 25,
			MaxIdleConns: 25,
			MaxIdleTime:  "15m",
		},
		Env:         "testing",
		ApiURL:      "http://localhost:8080",
		FrontendURL: "http://localhost:3000",
		Mail: MailConfig{
			Exp: 24 * time.Hour,
		},
		Auth: AuthConfig{
			Token: TokenConfig{
				Secret: "testsecret",
				Exp:    24 * time.Hour,
				Iss:    "squares-api",
			},
		},
	}

	store := store.NewStorage(postgres)

	svc := services.NewServices(&store)
	authenticator := auth.NewJWTAuthenticator(cfg.Auth.Token.Secret, cfg.Auth.Token.Iss, cfg.Auth.Token.Iss)

	testApp := NewApplication(cfg, store, svc, logger, authenticator)

	return testApp
}

func initTestDBTables(db *sql.DB) error {
	if _, err := db.Exec("CREATE EXTENSION IF NOT EXISTS citext"); err != nil {
		return fmt.Errorf("failed to create citext extension: %w", err)
	}
	schema := `
	-- users table
	CREATE TABLE IF NOT EXISTS users (
		id bigserial PRIMARY KEY,
		email citext UNIQUE NOT NULL,
		username varchar(255) UNIQUE NOT NULL,
		display_name varchar(255),
		password bytea NOT NULL,
		created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
		is_active BOOLEAN NOT NULL DEFAULT FALSE
	);

	-- user_invitations table for authentication/activation
	CREATE TABLE IF NOT EXISTS user_invitations (
		token bytea PRIMARY KEY,
		user_id bigint NOT NULL,
		expiry timestamp(0) with time zone NOT NULL
	);

	-- tasks table
	CREATE TABLE IF NOT EXISTS tasks (
		id bigserial PRIMARY KEY,
		user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		title varchar(255) NOT NULL,
		is_optional boolean DEFAULT FALSE,
		created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
		updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
		completions_count integer NOT NULL DEFAULT 0,
		completed BOOLEAN NOT NULL DEFAULT FALSE,
		tier smallint NOT NULL DEFAULT 2
	);
`

	_, err := db.Exec(schema)
	return err
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func newInMemTestDatabase(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(5433).
		Username("postgres").
		Password("postgres").
		Database("testdb"))

	err := postgres.Start()
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL: %v", err)
	}

	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=postgres dbname=testdb sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	err = initTestDBTables(db)
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	return db, func() {
		db.Close()
		postgres.Stop()
	}
}

func createTestUser(t *testing.T, db *sql.DB, app *Application) int64 {
	t.Helper()

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
		if err := app.store.Users.Create(ctx, tx, user); err != nil {
			return err
		}
		userID = user.ID
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return userID
}

func getAuthToken(t *testing.T, app *Application, mux http.Handler) string {
	t.Helper()

	loginBody := fmt.Sprintf(`{"email":"testuser@example.com","password":"testpassword"}`)
	req, _ := http.NewRequest("POST", "/v1/auth/token", strings.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("Login failed: %d - %s", rr.Code, rr.Body.String())
	}

	var tokenResp struct {
		Data string `json:"data"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&tokenResp); err != nil {
		t.Fatalf("Failed to decode token: %v", err)
	}
	return tokenResp.Data
}

func createTestTask(t *testing.T, ctx context.Context, app *Application, userID int64) *store.Task {
	t.Helper()

	task := &store.Task{
		Title:      "Test Task",
		IsOptional: false,
		Tier:       1,
		UserID:     userID,
	}

	if err := app.store.Tasks.CreateTasks(ctx, []*store.Task{task}); err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	return task
}

func checkResponseCode(t *testing.T, expected, actual int) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}
