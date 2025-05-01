package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/muzhiknastya/squares-my-beloved/internal/db"
	"github.com/muzhiknastya/squares-my-beloved/internal/env"
	"github.com/muzhiknastya/squares-my-beloved/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Failed to load .env.", err, "Default values are used.")
	}
}

//	@title						SquaresMyBeloved API
//	@description				API for SquaresMyBeloved, a web app to track productivity
//	@termsOfService				http://swagger.io/terms/
//	@contact.name				API Support
//	@contact.url				http://www.swagger.io/support
//	@contact.email				support@swagger.io
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//
//	@BasePath					/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description
func main() {
	// Config
	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/squares?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp: time.Hour * 24 * 3,
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	logger.Info("database connection pool has established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
