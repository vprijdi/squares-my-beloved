package api

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
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
//	@Router			/auth/register [post]
func (app *Application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
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

	err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.Config.Mail.Exp)

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
//
//	@Summary		Activate user account
//	@Description	Activates a user account using the activation token
//	@Tags			auth
//	@Param			token	path	string	true	"Activation token"
//	@Success		204		"Account activated successfully"
//	@Failure		400		{object}	map[string]string	"Invalid token"
//	@Failure		404		{object}	map[string]string	"Token not found"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/auth/activate/{token} [get]
func (app *Application) activateUserHandler(w http.ResponseWriter, r *http.Request) {

	token := chi.URLParam(r, "token")

	err := app.services.UserServices.ActivateAndSetup(r.Context(), token)

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

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// CreateToken godoc
//
//	@Summary		Create authentication token
//	@Description	Generates a JWT token for authenticated users
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"Credentials"
//
//	@Success		201		{string}	string					"Token"
//
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/auth/token [post]
func (app *Application) createTokenHandler(w http.ResponseWriter, r *http.Request) {

	var payload CreateUserTokenPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// fetch the user (check if user exist) from payload
	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)

	if err != nil {

		switch err {
		case store.ErrNotFound:
			app.unauthorizedErrorResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := user.Password.Compare(payload.Password); err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	// generate the token
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.Config.Auth.Token.Exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.Config.Auth.Token.Iss,
		"aud": app.Config.Auth.Token.Iss,
	}

	token, err := app.authenticator.GenerateToken(claims)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)
	}
}
