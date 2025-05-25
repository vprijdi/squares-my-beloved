package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/muzhiknastya/squares-my-beloved/internal/store"
)

type contextKey string

const (
	taskCtxKey contextKey = "task"
	userCtxKey contextKey = "user"
	goalCtxKey contextKey = "goal"
)

func (app *Application) taskContextMiddleware(next http.Handler) http.Handler {
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

		ctx = context.WithValue(ctx, taskCtxKey, task)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *Application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userID")
		userID, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
		}

		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *Application) authTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		token := parts[1]
		jwtToken, err := app.authenticator.ValidateToken(token)

		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)

		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, userID)

		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *Application) goalContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		goalID, err := strconv.ParseInt(chi.URLParam(r, "goalID"), 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		goal, err := app.store.Goals.GetByID(r.Context(), goalID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		user := getUserFromCtx(r)
		if goal.UserID != user.ID {
			app.notFoundResponse(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), goalCtxKey, goal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtxKey).(*store.User)
	return user
}

func getTaskFromCtx(r *http.Request) *store.Task {
	task, _ := r.Context().Value(taskCtxKey).(*store.Task)
	return task
}

func getGoalFromCtx(r *http.Request) *store.Goal {
	return r.Context().Value(goalCtxKey).(*store.Goal)
}
