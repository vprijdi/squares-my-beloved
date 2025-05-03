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

// func NewInMemTestDatabase(ctx context.Context, cfg *)
// // Global vars for testing
// var (
// 	postgres *embeddedpostgres.EmbeddedPostgres
// 	testDB   *sql.DB
// 	testApp  *Application
// 	testMux  http.Handler
// )

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	// Initialize database schema for testing

// 	// Initialize the store, services, and authenticator
// 	store := store.NewStorage(testDB)

// 	// Create test configuration
// 	cfg := Config{
// 		Addr: "localhost:8080", // Not actually used for tests
// 		DB: DBConfig{
// 			Addr:         "localhost:5433",
// 			MaxOpenConns: 25,
// 			MaxIdleConns: 25,
// 			MaxIdleTime:  "15m",
// 		},
// 		Env:         "testing",
// 		ApiURL:      "http://localhost:8080",
// 		FrontendURL: "http://localhost:3000",
// 		Mail: MailConfig{
// 			Exp: 24 * time.Hour,
// 		},
// 		Auth: AuthConfig{
// 			Token: TokenConfig{
// 				Secret: "testsecret",
// 				Exp:    24 * time.Hour,
// 				Iss:    "squares-api",
// 			},
// 		},
// 	}

// 	// Initialize the authenticator
// 	authenticator := auth.NewJWTAuthenticator(
// 		cfg.Auth.Token.Secret,
// 		cfg.Auth.Token.Exp.String(),
// 		cfg.Auth.Token.Iss,
// 	)

// 	// Initialize the services layer
// 	svc := services.NewServices(&store)

// 	// Initialize the application

// 	// Mount routes for testing
// 	fmt.Println("Mounting routes...")
// 	testMux = testApp.Mount()
// 	fmt.Printf("Router initialized: %v\n", testMux != nil)

// 	// Run the tests
// 	exitCode := m.Run()

// 	// Cleanup phase
// 	testDB.Close()
// 	postgres.Stop()

// 	os.Exit(exitCode)
// }

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

// // Helper function to create a test user and return authentication token
// func createTestUser(t *testing.T) (map[string]interface{}, string) {
// 	// Register a user
// 	userData := map[string]interface{}{
// 		"name":     "Test User",
// 		"email":    "test@example.com",
// 		"password": "password123",
// 	}

// 	userJSON, _ := json.Marshal(userData)
// 	req, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(userJSON))
// 	req.Header.Set("Content-Type", "application/json")

// 	statusCode, resp := executeRequest(t, req)
// 	require.Equal(t, http.StatusCreated, statusCode)

// 	// Activate the user directly in the database
// 	_, err := testDB.Exec("UPDATE users SET activated = true WHERE email = $1", userData["email"])
// 	require.NoError(t, err)

// 	// Get authentication token
// 	loginData := map[string]interface{}{
// 		"email":    userData["email"],
// 		"password": userData["password"],
// 	}

// 	loginJSON, _ := json.Marshal(loginData)
// 	req, _ = http.NewRequest("POST", "/v1/auth/token", bytes.NewBuffer(loginJSON))
// 	req.Header.Set("Content-Type", "application/json")

// 	statusCode, resp = executeRequest(t, req)
// 	require.Equal(t, http.StatusCreated, statusCode)
// 	require.Contains(t, resp, "token")

//		token := resp["token"].(string)
//		return userData, token
//	}
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
