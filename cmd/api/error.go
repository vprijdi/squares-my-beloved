package api

import (
	"errors"
	"net/http"
)

var (
	ErrEmptyTaskList = errors.New("at least one task is required")
)

func (app *Application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err,
	)

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("bad request error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err,
	)

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("not found error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err,
	)

	writeJSONError(w, http.StatusNotFound, "resource not found")
}

func (app *Application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}
