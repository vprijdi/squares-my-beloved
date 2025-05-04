package api

import (
	"net/http"

	"github.com/muzhiknastya/squares-my-beloved/internal/store"
)

// GetUser godoc
//
//	@Summary		Get user details
//	@Description	Get a user's information by ID
//	@Tags			users
//	@Produce		json
//	@Param			userID	path		string				true	"User ID"
//	@Success		200		{object}	store.User			"User details"
//	@Failure		404		{object}	map[string]string	"User not found"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/users/{userID} [get]
func (app *Application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
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
