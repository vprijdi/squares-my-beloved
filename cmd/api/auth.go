package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/muzhiknastya/squares-my-beloved/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=72"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// RegisterUser godoc
//
//	@Summary		Register a new user
//	@Description	Creates a new user account and returns an activation token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"Registration data"
//	@Success		201		{object}	UserWithToken		"User created with activation token"
//	@Failure		400		{object}	map[string]string	"Invalid request/duplicate email or username"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/v1/auth/register [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	plainToken := uuid.New().String()

	// hash the token for storage but keep the plain token for email
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp)

	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestResponse(w, r, err)
		case store.ErrDuplicateUsername:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}
	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// ActivateUser godoc
//	@Summary		Activate user account
//	@Description	Activates a user account using the activation token
//	@Tags			auth
//	@Param			token	path	string	true	"Activation token"
//	@Success		204		"Account activated successfully"
//	@Failure		400		{object}	map[string]string	"Invalid token"
//	@Failure		404		{object}	map[string]string	"Token not found"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/v1/auth/activate/{token} [get]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {

	token := chi.URLParam(r, "token")

	err := app.store.Users.Activate(r.Context(), token)

	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
		app.internalServerError(w, r, err)
	}
}
