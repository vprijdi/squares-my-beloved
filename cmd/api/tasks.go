package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/muzhiknastya/squares-my-beloved/internal/store"
)

var (
	ErrEmptyTaskList = errors.New("at least one task is required")
)

type createTaskPayload struct {
	Title      string `json:"title" validate:"required,min=1,max=160"`
	IsOptional bool   `json:"is_optional"`
}

func (app *application) createBatchTasksHandler(w http.ResponseWriter, r *http.Request) {
	var payloads []createTaskPayload
	if err := readJSONArray(w, r, &payloads); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if len(payloads) == 0 {
		app.badRequestResponse(w, r, ErrEmptyTaskList)
		return
	}

	tasks := make([]*store.Task, len(payloads))
	for i, payload := range payloads {
		if err := Validate.Struct(payload); err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		tasks[i] = &store.Task{
			Title:      payload.Title,
			IsOptional: payload.IsOptional,
			UserID:     1, // TODO: change when auth
		}
	}

	ctx := r.Context()

	err := app.store.Tasks.CreateTasks(ctx, tasks)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"tasks": tasks,
		"count": len(tasks),
	}); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	task := getTaskFromCtx(r)

	if err := app.jsonResponse(w, http.StatusCreated, task); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) tasksContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "taskID")
		taskID, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
		}

		ctx := r.Context()

		task, err := app.store.Tasks.GetByID(ctx, taskID)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, "task", task)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getTaskFromCtx(r *http.Request) *store.Task {
	task, _ := r.Context().Value("task").(*store.Task)
	return task
}
