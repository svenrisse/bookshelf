package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/markbates/goth/gothic"
	"github.com/svenrisse/bookshelf/internal/models"
)

func (app *application) AuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	provider := params.ByName("provider")

	ctx := context.WithValue(context.Background(), "provider", provider)
	r = r.WithContext(ctx)

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		app.logger.Error(err.Error())
		return
	}

	id, err := strconv.Atoi(user.UserID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	exists, err := app.models.Users.Exists(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !exists {
		err = app.models.Users.Insert(&models.User{
			ID:       id,
			Name:     user.NickName,
			Avatar:   user.AvatarURL,
			Provider: user.Provider,
		})
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	http.Redirect(w, r, "https://bookshelf.svenrisse.com/v1/healthcheck", http.StatusTemporaryRedirect)
}

func (app *application) AuthLogout(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	provider := params.ByName("provider")

	ctx := context.WithValue(context.Background(), "provider", provider)
	r = r.WithContext(ctx)
	gothic.Logout(w, r)
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (app *application) Auth(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	provider := params.ByName("provider")

	ctx := context.WithValue(context.Background(), "provider", provider)
	r = r.WithContext(ctx)

	if _, err := gothic.CompleteUserAuth(w, r); err == nil {
		http.Redirect(w, r, "http://localhost:4000/", http.StatusTemporaryRedirect)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}
