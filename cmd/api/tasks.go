package api

import (
	"errors"
	"net/http"

	"github.com/muzhiknastya/squares-my-beloved/internal/store"
)

type createTaskPayload struct {
	Title      string `json:"title" validate:"required,min=1,max=160"`
	IsOptional bool   `json:"is_optional"`
	Tier       int32  `json:"tier" validate:"oneof=0 1 2"`
}

// CreateBatchTasks godoc
//
//	@Summary		Create multiple tasks in a batch
//	@Description	Create several tasks with a single request
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			tasks	body		[]createTaskPayload		true	"List of tasks to create"
//	@Success		201		{object}	map[string]interface{}	"Created tasks with count"
//	@Failure		400		{object}	map[string]string		"Invalid request"
//	@Failure		401		{object}	map[string]string		"Unauthorized"
//	@Failure		500		{object}	map[string]string		"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/tasks [post]
func (app *Application) createBatchTasksHandler(w http.ResponseWriter, r *http.Request) {
	var payloads []createTaskPayload
	if err := readJSONArray(w, r, &payloads); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if len(payloads) == 0 {
		app.badRequestResponse(w, r, ErrEmptyTaskList)
		return
	}

	user := getUserFromCtx(r)
	tasks := make([]*store.Task, len(payloads))
	for i, payload := range payloads {
		if err := Validate.Struct(payload); err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		tasks[i] = &store.Task{
			Title:      payload.Title,
			IsOptional: payload.IsOptional,
			Tier:       payload.Tier,
			UserID:     user.ID,
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

// GetTask godoc
//
//	@Summary		Get a specific task
//	@Description	Get task by ID
//	@Tags			tasks
//	@Produce		json
//	@Param			id	path		int					true	"Task ID"
//	@Success		200	{object}	store.Task			"The requested task"
//	@Failure		404	{object}	map[string]string	"Task not found"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/tasks/{id} [get]
func (app *Application) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	task := getTaskFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, task); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// CompleteTask godoc
//
//	@Summary		Complete a task
//	@Description	Mark a task as completed
//	@Tags			tasks
//	@Param			id	path	int	true	"Task ID"
//	@Success		204	"No content"
//	@Failure		404	{object}	map[string]string	"Task not found"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/tasks/{id}/complete [patch]
func (app *Application) completeTaskHandler(w http.ResponseWriter, r *http.Request) {
	task := getTaskFromCtx(r)

	ctx := r.Context()

	if err := app.store.Tasks.Complete(ctx, task.ID); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// listUserTasksHandler godoc
//
//	@Summary		List tasks for a specific user
//	@Description	Get paginated and filtered list of tasks for the specified user ID
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			userID		path	string	true	"User ID"
//	@Param			limit		query	int		false	"Results limit (1-100)"	default(100)	minimum(1)	maximum(100)
//	@Param			offset		query	int		false	"Pagination offset"		default(0)		minimum(0)
//	@Param			date		query	string	false	"Date filter (today/yesterday/YYYY-MM-DD)"
//	@Param			completed	query	boolean	false	"Filter by completion status"
//	@Security		ApiKeyAuth
//	@Success		201	"Tasks returned successfully"
//	@Failure		400	"Invalid request parameters"
//	@Failure		500	"Server error"
//	@Router			/users/{userID}/tasks [get]
func (app *Application) listUserTasksHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	tf := store.TaskFilters{
		Limit:  100,
		Offset: 0,
	}

	tf, err := tf.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(tf); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	filteredResult, err := app.store.Tasks.GetAllUserTasks(ctx, user.ID, &tf)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, filteredResult); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
