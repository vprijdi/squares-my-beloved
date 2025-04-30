package main

import (
	"net/http"

	"github.com/muzhiknastya/squares-my-beloved/internal/store"
)

type CreateTaskPayload struct {
	Title      string `json:"title"`
	IsOptional bool   `json:"is_optional"`
}

func (app *application) createTaskhandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateTaskPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	task := &store.Task{
		Title:      payload.Title,
		IsOptional: payload.IsOptional,
		// TODO: change after auth
		UserID: 1,
	}

	ctx := r.Context()

	if err := app.store.Tasks.Create(ctx, task); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, task); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
