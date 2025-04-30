package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/muzhiknastya/squares-my-beloved/internal/store"
)

type createTaskPayload struct {
	Title      string `json:"title"`
	IsOptional bool   `json:"is_optional"`
}

func (app *application) createTaskhandler(w http.ResponseWriter, r *http.Request) {
	var payload createTaskPayload
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

	if err := app.jsonResponse(w, http.StatusCreated, task); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "taskID")
	taskID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	ctx := r.Context()

	task, err := app.store.Tasks.GetByID(ctx, taskID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, task); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
