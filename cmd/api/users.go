package api

import (
	"net/http"
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
