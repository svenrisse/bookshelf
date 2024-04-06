package main

import (
	"context"
	"net/http"

	"github.com/svenrisse/bookshelf/internal/data"
)

type contextKey string

const UserContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(UserContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
