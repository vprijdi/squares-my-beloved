package api

import (
	"errors"
	"net/http"

	"github.com/muzhiknastya/squares-my-beloved/internal/store"
)

type createGoalPayload struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
}

// CreateBatchGoals godoc
//
//	@Summary		Create multiple goals in a batch
//	@Description	Create several goals with a single request
//	@Tags			goals
//	@Accept			json
//	@Produce		json
//	@Param			goals	body		[]createGoalPayload		true	"List of goals to create"
//	@Success		201		{object}	map[string]interface{}	"Created goals with count"
//	@Failure		400		{object}	map[string]string		"Invalid request"
//	@Failure		401		{object}	map[string]string		"Unauthorized"
//	@Failure		500		{object}	map[string]string		"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/goals [post]
func (app *Application) createBatchGoalsHandler(w http.ResponseWriter, r *http.Request) {
	var payloads []createGoalPayload
	if err := readJSONArray(w, r, &payloads); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if len(payloads) == 0 {
		app.badRequestResponse(w, r, ErrEmptyGoalList)
		return
	}

	user := getUserFromCtx(r)
	goals := make([]*store.Goal, len(payloads))
	for i, payload := range payloads {
		if err := Validate.Struct(payload); err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		goals[i] = &store.Goal{
			Title:  payload.Title,
			UserID: user.ID,
		}
	}

	ctx := r.Context()

	err := app.store.Goals.CreateGoals(ctx, goals)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"goals": goals,
		"count": len(goals),
	}); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetGoal godoc
//
//	@Summary		Get a specific goal
//	@Description	Get goal by ID
//	@Tags			goals
//	@Produce		json
//	@Param			id	path		int					true	"Goal ID"
//	@Success		200	{object}	store.Goal			"The requested goal"
//	@Failure		404	{object}	map[string]string	"Goal not found"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/goals/{id} [get]
func (app *Application) getGoalHandler(w http.ResponseWriter, r *http.Request) {
	goal := getGoalFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, goal); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// AchieveGoal godoc
//
//	@Summary		Mark goal as achieved
//	@Description	Mark a goal as achieved
//	@Tags			goals
//	@Param			id	path	int	true	"Goal ID"
//	@Success		204	"No content"
//	@Failure		404	{object}	map[string]string	"Goal not found"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/goals/{id}/achieve [patch]
func (app *Application) achieveGoalHandler(w http.ResponseWriter, r *http.Request) {
	goal := getGoalFromCtx(r)

	ctx := r.Context()

	if err := app.store.Goals.Achieve(ctx, goal.ID); err != nil {
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

// listUserGoalsHandler godoc
//
//	@Summary		List goals for a specific user
//	@Description	Get paginated and filtered list of goals for the specified user ID
//	@Tags			goals
//	@Accept			json
//	@Produce		json
//	@Param			userID		path	string	true	"User ID"
//	@Param			limit		query	int		false	"Results limit (1-100)"	default(100)	minimum(1)	maximum(100)
//	@Param			offset		query	int		false	"Pagination offset"		default(0)		minimum(0)
//	@Param			achieved	query	boolean	false	"Filter by achievement status"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	[]store.Goal		"Goals returned successfully"
//	@Failure		400	{object}	map[string]string	"Invalid request parameters"
//	@Failure		401	{object}	map[string]string	"Unauthorized"
//	@Failure		500	{object}	map[string]string	"Server error"
//	@Router			/users/{userID}/goals [get]
func (app *Application) listUserGoalsHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	gf := store.GoalFilters{
		Limit:  100,
		Offset: 0,
	}

	gf, err := gf.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(gf); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	goals, err := app.store.Goals.GetAllUserGoals(ctx, user.ID, &gf)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, goals); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
