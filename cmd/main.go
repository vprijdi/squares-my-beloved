package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/muzhiknastya/squares-my-beloved/cmd/api"
	"github.com/muzhiknastya/squares-my-beloved/internal/auth"
	"github.com/muzhiknastya/squares-my-beloved/internal/db"
	"github.com/muzhiknastya/squares-my-beloved/internal/env"
	"github.com/muzhiknastya/squares-my-beloved/internal/services"
	"github.com/muzhiknastya/squares-my-beloved/internal/store"
	"go.uber.org/zap"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Failed to load .env.", err, "Default values are used.")
	}
}

// @title						SquaresMyBeloved API
// @description				API for SquaresMyBeloved, a web app to track productivity
// @termsOfService				http://swagger.io/terms/
// @contact.name				API Support
// @contact.url				http://www.swagger.io/support
// @contact.email				support@swagger.io
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//
// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description

func main() {
	// Config
	cfg := api.Config{
		Addr:        env.GetString("ADDR", ":8080"),
		ApiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		FrontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		DB: api.DBConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/squares?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		Env: env.GetString("ENV", "development"),
		Mail: api.MailConfig{
			Exp: time.Hour * 24 * 3,
		},
		Auth: api.AuthConfig{
			Token: api.TokenConfig{
				Secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
				Exp:    time.Hour * 24 * 3,
				Iss:    "gophersocial",
			},
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	database, err := db.New(cfg.DB.Addr, cfg.DB.MaxOpenConns, cfg.DB.MaxIdleConns, cfg.DB.MaxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer database.Close()
	logger.Info("database connection pool has established")

	storage := store.NewStorage(database)
	services := services.NewServices(&storage, database)
	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.Auth.Token.Secret, cfg.Auth.Token.Iss, cfg.Auth.Token.Iss)

	app := api.NewApplication(cfg, storage, services, logger, jwtAuthenticator)

	mux := app.Mount()
	logger.Fatal(app.Run(mux))
}
