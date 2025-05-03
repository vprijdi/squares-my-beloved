package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/muzhiknastya/squares-my-beloved/docs"
	"github.com/muzhiknastya/squares-my-beloved/internal/auth"
	"github.com/muzhiknastya/squares-my-beloved/internal/services"
	"github.com/muzhiknastya/squares-my-beloved/internal/store"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

const version = "0.0.1"

type Application struct {
	Config        Config
	store         store.Storage
	services      services.Services
	logger        *zap.SugaredLogger
	authenticator auth.Authenticator
}

type Config struct {
	Addr        string
	DB          DBConfig
	Env         string
	ApiURL      string
	FrontendURL string
	Mail        MailConfig
	Auth        AuthConfig
}

type AuthConfig struct {
	Token TokenConfig
}

type TokenConfig struct {
	Secret string
	Exp    time.Duration
	Iss    string
}

type MailConfig struct {
	Exp time.Duration
}

type DBConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

// NewApplication creates a new Application instance
func NewApplication(config Config, store store.Storage, services services.Services, logger *zap.SugaredLogger, authenticator auth.Authenticator) *Application {
	return &Application{
		Config:        config,
		store:         store,
		services:      services,
		logger:        logger,
		authenticator: authenticator,
	}
}

func (app *Application) Mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.Config.Addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.registerUserHandler)
			r.Get("/activate/{token}", app.activateUserHandler)
			r.Post("/token", app.createTokenHandler)
		})

		r.Route("/tasks", func(r chi.Router) {
			r.Use(app.authTokenMiddleware)
			r.Post("/", app.createBatchTasksHandler)

			r.Route("/{taskID}", func(r chi.Router) {
				r.Use(app.tasksContextMiddleware)

				r.Get("/", app.getTaskHandler)
				r.Patch("/complete", app.completeTaskHandler)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)
				r.Get("/tasks", app.listUserTasksHandler)
				r.Get("/", app.getUserHandler)
			})
		})
	})

	return r
}

func (app *Application) Run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.Config.ApiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	app.logger.Infow("server has started",
		"addr", app.Config.Addr,
		"env", app.Config.Env,
	)

	return srv.ListenAndServe()
}
