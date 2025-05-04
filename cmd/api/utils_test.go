package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
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
