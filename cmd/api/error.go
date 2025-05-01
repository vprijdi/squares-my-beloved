package main

import (
	"errors"
	"net/http"
)

var (
	ErrEmptyTaskList = errors.New("at least one task is required")
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err,
	)

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("bad request error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err,
	)

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("not found error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err,
	)

	writeJSONError(w, http.StatusNotFound, "resource not found")
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}
