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

func main() {
	// Config
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
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
