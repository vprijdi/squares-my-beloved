package main

import "net/http"

//	@Summary		Health check endpoint
//	@Description	Returns the current health status of the application along with environment and version information
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string	"Returns application health status"
//	@Failure		500	{object}	map[string]string	"When there's a server error while generating the response"
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
